const state = {
    allFiles: [],         // 存储所有文件数据
    viewFiles: [],        // 存储当前过滤或搜索后的文件数据
    selectedIds: new Set(), // 存储选中的文件ID
    filterType: 'all',    // 当前文件类型过滤设置
    pagination: {         // 分页状态
        page: 1,
        limit: 10
    },
    currentDetailId: null,      // 当前详情弹窗显示的文件ID
    currentDetailFilePath: null // 当前详情弹窗显示的文件路径
};

const $ = (selector) => document.querySelector(selector);
const $all = (selector) => document.querySelectorAll(selector);

document.addEventListener('DOMContentLoaded', () => {
    fetchFiles();
    setupEvents();
});

/**
 * 设置页面事件监听器，包括搜索、过滤、文件上传、全选/取消全选和详情弹窗下载按钮。
 */
function setupEvents() {
    $('#search-input').addEventListener('input', (e) => {
        state.pagination.page = 1;
        applyFilters(e.target.value);
    });

    $all('.tab').forEach(btn => {
        btn.addEventListener('click', (e) => {
            $all('.tab').forEach(t => t.classList.remove('active'));
            e.currentTarget.classList.add('active');
            state.filterType = e.currentTarget.dataset.type;
            state.pagination.page = 1;
            applyFilters($('#search-input').value);
        });
    });

    const dropZone = $('#drop-zone');
    const fileInput = $('#file-input');
    dropZone.addEventListener('click', () => fileInput.click());
    fileInput.addEventListener('change', (e) => handleUpload(e.target.files));
    ['dragenter', 'dragover'].forEach(e => dropZone.addEventListener(e, (ev) => { ev.preventDefault(); dropZone.style.borderColor = '#6366f1'; }));
    ['dragleave', 'drop'].forEach(e => dropZone.addEventListener(e, (ev) => { ev.preventDefault(); dropZone.style.borderColor = ''; }));
    dropZone.addEventListener('drop', (e) => handleUpload(e.dataTransfer.files));

    $('#select-all').addEventListener('change', (e) => {
        const paginatedData = getPaginatedData();
        if(e.target.checked) paginatedData.forEach(f => state.selectedIds.add(f.id));
        else paginatedData.forEach(f => state.selectedIds.delete(f.id));
        renderTable();
        updateBatchBar();
    });

    $('#modal-download-btn').addEventListener('click', () => {
        if(state.currentDetailId) triggerDownload(state.currentDetailId);
    });
}

/**
 * 显示自定义确认弹窗。
 * @param {string} title - 弹窗标题。
 * @param {string} message - 弹窗消息。
 * @returns {Promise<boolean>} - 用户点击确认返回 true，点击取消返回 false。
 */
function showConfirmModal(title, message) {
    return new Promise((resolve) => {
        const modal = $('#confirm-modal');
        $('#confirm-title').innerText = title;
        $('#confirm-msg').innerText = message;
        modal.style.display = 'flex';

        const cleanup = () => {
            $('#btn-cancel-confirm').onclick = null;
            $('#btn-ok-confirm').onclick = null;
            modal.style.display = 'none';
        };

        $('#btn-cancel-confirm').onclick = () => { cleanup(); resolve(false); };
        $('#btn-ok-confirm').onclick = () => { cleanup(); resolve(true); };
    });
}

/**
 * 从服务器获取文件列表。
 */
window.fetchFiles = async () => {
    try {
        const btnIcon = $('.btn-secondary-pill i');
        if(btnIcon) btnIcon.classList.add('fa-spin');

        const res = await fetch('/api/v1/files?page=1&limit=99999');
        const data = await res.json();

        if (data.success) {
            state.allFiles = data.files.sort((a,b) => new Date(b.created_at) - new Date(a.created_at));
            applyFilters($('#search-input').value);
            updateStats();
        }
    } catch (e) {
        showToast('无法连接服务器或加载文件', 'error');
        console.error("Error fetching files:", e);
    } finally {
        const btnIcon = $('.btn-secondary-pill i');
        if(btnIcon) setTimeout(() => btnIcon.classList.remove('fa-spin'), 500);
    }
};

/**
 * 根据关键字和文件类型过滤文件列表。
 * @param {string} keyword - 搜索关键字。
 */
function applyFilters(keyword = '') {
    const key = keyword.toLowerCase();

    state.viewFiles = state.allFiles.filter(f => {
        let typeMatch = true;
        if (state.filterType === 'pdf') typeMatch = f.extension.toLowerCase() === 'pdf';
        if (state.filterType === 'word') typeMatch = ['doc', 'docx'].includes(f.extension.toLowerCase());
        const nameMatch = f.original_name.toLowerCase().includes(key);
        return typeMatch && nameMatch;
    });

    state.selectedIds.clear();
    renderTable();
    renderPagination();
    updateBatchBar();
}

/**
 * 处理文件上传。
 * @param {FileList} files - 要上传的文件列表。
 */
async function handleUpload(files) {
    if(!files.length) return;
    $('#upload-mask').style.display = 'flex';

    const formData = new FormData();
    Array.from(files).forEach(f => formData.append('files', f));

    try {
        const res = await fetch('/api/v1/upload/batch', { method: 'POST', body: formData });
        const data = await res.json();
        if(data.success || (data.files && data.files.length)) {
            showToast(`成功上传 ${data.files.length} 个文件`);
            fetchFiles();
        } else {
            showToast('上传失败：' + (data.message || '未知错误'), 'error');
        }
    } catch (e) {
        showToast('网络错误或服务器无响应', 'error');
        console.error("Upload error:", e);
    } finally {
        $('#upload-mask').style.display = 'none';
        $('#file-input').value = '';
    }
}

/**
 * 获取当前页的经过过滤的文件数据。
 * @returns {Array} - 当前页的文件数据。
 */
function getPaginatedData() {
    const start = (state.pagination.page - 1) * state.pagination.limit;
    return state.viewFiles.slice(start, start + state.pagination.limit);
}

/**
 * 渲染文件列表表格。
 */
function renderTable() {
    const list = getPaginatedData();
    const tbody = $('#file-list-body');
    tbody.innerHTML = '';

    if(!list.length) {
        tbody.innerHTML = `<tr><td colspan="6" class="text-center" style="padding:40px;color:#9ca3af;">没有找到匹配的文件</td></tr>`;
        $('#select-all').checked = false;
        return;
    }

    list.forEach(file => {
        const isSelected = state.selectedIds.has(file.id);
        const ext = file.extension.toLowerCase();
        const iconInfo = ext === 'pdf' ?
            { iconClass: 'fa-file-pdf', badgeClass: 'fb-pdf' } :
            { iconClass: 'fa-file-word', badgeClass: 'fb-word' };

        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td class="text-center"><input type="checkbox" class="row-check" value="${file.id}" ${isSelected?'checked':''}></td>
            <td>
                <div style="display:flex;align-items:center;gap:12px;">
                    <div class="file-badge ${iconInfo.badgeClass}"><i class="fa-solid ${iconInfo.iconClass}"></i></div>
                    <span style="font-weight:500;">${file.original_name}</span>
                </div>
            </td>
            <td><span style="color:#6b7280; font-size:12px; text-transform:uppercase;">${ext}</span></td>
            <td><span style="color:#6b7280; font-size:13px;">${formatSize(file.size)}</span></td>
            <td><span style="color:#6b7280; font-size:13px;">${formatDate(file.created_at)}</span></td>
            <td class="text-right">
                <button class="action-btn btn-view" onclick="openDetail('${file.id}')" title="详情"><i class="fa-solid fa-eye"></i></button>
                <button class="action-btn btn-down" onclick="triggerDownload('${file.id}')" title="下载"><i class="fa-solid fa-download"></i></button>
                <button class="action-btn btn-del" onclick="deleteFile('${file.id}')" title="删除"><i class="fa-solid fa-trash-can"></i></button>
            </td>
        `;

        tr.querySelector('.row-check').addEventListener('change', (e) => {
            if(e.target.checked) state.selectedIds.add(file.id);
            else state.selectedIds.delete(file.id);
            updateBatchBar();
        });
        tbody.appendChild(tr);
    });

    $('#select-all').checked = list.length > 0 && list.every(f => state.selectedIds.has(f.id));
}

/**
 * 渲染分页控件。
 */
function renderPagination() {
    const totalFilesCount = state.viewFiles.length;
    const itemsPerPage = state.pagination.limit;
    const totalPages = Math.ceil(totalFilesCount / itemsPerPage);
    const currentPage = state.pagination.page;
    const paginationContainer = $('#pagination');

    paginationContainer.innerHTML = '';

    if(totalPages <= 1 && totalFilesCount <= itemsPerPage) {
        return;
    }

    paginationContainer.appendChild(createPageBtn('<i class="fa-solid fa-chevron-left"></i>', currentPage > 1, () => changePage(currentPage - 1)));

    const startPage = Math.max(1, currentPage - 1);
    const endPage = Math.min(totalPages, currentPage + 1);

    for (let i = startPage; i <= endPage; i++) {
        const btn = createPageBtn(i, true, () => changePage(i));
        if (i === currentPage) btn.classList.add('active');
        paginationContainer.appendChild(btn);
    }

    paginationContainer.appendChild(createPageBtn('<i class="fa-solid fa-chevron-right"></i>', currentPage < totalPages, () => changePage(currentPage + 1)));
}

/**
 * 创建分页按钮。
 * @param {string|number} htmlContent - 按钮显示的 HTML 内容或文本。
 * @param {boolean} isEnabled - 按钮是否可用。
 * @param {Function} onClickCallback - 按钮点击回调函数。
 * @returns {HTMLButtonElement} - 创建的按钮元素。
 */
function createPageBtn(htmlContent, isEnabled, onClickCallback) {
    const button = document.createElement('button');
    button.className = 'page-btn';
    button.innerHTML = htmlContent;
    button.disabled = !isEnabled;
    if(isEnabled) button.onclick = onClickCallback;
    return button;
}

/**
 * 切换到指定页码并重新渲染表格和分页。
 * @param {number} newPage - 新的页码。
 */
function changePage(newPage) {
    state.pagination.page = newPage;
    renderTable();
    renderPagination();
}

/**
 * 打开文件详情弹窗。
 * @param {string} id - 文件ID。
 */
window.openDetail = (id) => {
    const file = state.allFiles.find(x => x.id === id);
    if (!file) return;

    state.currentDetailId = id;
    state.currentDetailFilePath = file.storage_path; // 这个应该已经包含了年月目录

    $('#modal-tag-ext').innerText = file.extension.toUpperCase();

    const detailContentHtml = `
        <div class="detail-key">文件ID</div><div class="detail-val">${file.id}</div>
        <div class="detail-key">名称</div><div class="detail-val">${file.original_name}</div>
        <div class="detail-key">类型</div><div class="detail-val">${file.mime_type}</div>
        <div class="detail-key">大小</div><div class="detail-val">${formatSize(file.size)}</div>
        <div class="detail-key">存储</div><div class="detail-val">${file.storage_path}</div>
        <div class="detail-key">时间</div><div class="detail-val">${formatDate(file.created_at)}</div>
    `;
    $('#detail-content').innerHTML = detailContentHtml;

    // 更新预览按钮的链接
    $('#modal-access-btn').onclick = () => {
        window.open(`/uploads/${file.storage_path}`, '_blank');
    };

    $('#detail-modal').style.display = 'flex';
};

/**
 * 关闭指定ID的模态框。
 * @param {string} id - 模态框的ID。
 */
window.closeModal = (id) => {
    $('#' + id).style.display = 'none';
};

/**
 * 触发文件下载。
 * @param {string} id - 文件ID。
 */
window.triggerDownload = (id) => {
    const anchor = document.createElement('a');
    anchor.href = `/api/v1/files/${id}/download`;
    anchor.click();
};

/**
 * 在新标签页中预览文件。
 */
window.viewFilePreview = () => {
    if (state.currentDetailFilePath) {
        // 使用相对路径
        window.open(`/uploads/${state.currentDetailFilePath}`, '_blank');
    }
};

/**
 * 删除文件。
 * @param {string} id - 文件ID。
 */
window.deleteFile = async (id) => {
    // 明确提示用户这是永久删除
    const confirmed = await showConfirmModal('确认删除?', '此文件将被<b>永久删除</b>，且无法恢复。您确定要继续吗？');
    if(!confirmed) return;
    try {
        const res = await fetch(`/api/v1/files/${id}`, { method: 'DELETE' });
        if (res.ok) {
            showToast('文件已成功删除');
            fetchFiles();
            state.selectedIds.delete(id);
            updateBatchBar();
        } else {
            const errorData = await res.json();
            showToast(`删除失败：${errorData.message || '未知错误'}`, 'error');
        }
    } catch(e) {
        showToast('删除操作发生网络错误', 'error');
        console.error("Delete file error:", e);
    }
};
/**
 * 确认并执行批量删除操作。
 */
window.confirmBatchDelete = async () => {
    const selectedCount = state.selectedIds.size;
    if (selectedCount === 0) {
        showToast('请先选择要删除的文件', 'error');
        return;
    }
    // 明确提示用户这是永久删除
    const confirmed = await showConfirmModal('确认批量删除?', `您选中的 ${selectedCount} 个文件将被<b>永久删除</b>，且无法恢复。您确定要继续吗？`);
    if(!confirmed) return;
    try {
        const deletePromises = Array.from(state.selectedIds).map(id =>
            fetch(`/api/v1/files/${id}`, { method: 'DELETE' })
                .then(res => res.ok ? {id: id, success: true} : res.json().then(err => ({id: id, success: false, message: err.message})))
        );
        const results = await Promise.all(deletePromises);
        let successCount = 0;
        let errorMessages = [];
        results.forEach(result => {
            if (result.success) {
                successCount++;
            } else {
                errorMessages.push(`文件ID ${result.id} 删除失败: ${result.message || '未知错误'}`);
            }
        });
        if (successCount > 0) {
            showToast(`已成功删除 ${successCount} 个文件`);
        }
        if (errorMessages.length > 0) {
            showToast(errorMessages.join('<br>'), 'error');
        }
        fetchFiles();
        state.selectedIds.clear();
        updateBatchBar();
    } catch (e) {
        showToast('批量删除操作发生网络错误', 'error');
        console.error("Batch delete error:", e);
    }
};

/**
 * 确认并执行批量删除操作。
 */
window.confirmBatchDelete = async () => {
    const selectedCount = state.selectedIds.size;
    if (selectedCount === 0) {
        showToast('请先选择要删除的文件', 'error');
        return;
    }

    const confirmed = await showConfirmModal('确认批量删除?', `您选中的 ${selectedCount} 个文件将被永久删除。`);
    if(!confirmed) return;

    try {
        for (const id of state.selectedIds) {
            await fetch(`/api/v1/files/${id}`, { method: 'DELETE' });
        }
        showToast(`已成功删除 ${selectedCount} 个文件`);
        fetchFiles();
        state.selectedIds.clear();
        updateBatchBar();
    } catch (e) {
        showToast('批量删除操作发生错误', 'error');
        console.error("Batch delete error:", e);
    }
};

/**
 * 更新页面上的文件统计数据。
 */
function updateStats() {
    $('#stat-total').innerText = state.allFiles.length;
    const totalSize = state.allFiles.reduce((accumulator, currentFile) => accumulator + currentFile.size, 0);
    $('#stat-size').innerText = (totalSize / (1024 * 1024)).toFixed(2) + ' MB';
}

/**
 * 更新批量操作栏的显示状态和选中数量。
 */
function updateBatchBar() {
    const count = state.selectedIds.size;
    $('#selected-count').innerText = count;
    $('#batch-bar').style.display = count > 0 ? 'inline-flex' : 'none';
}

/**
 * 显示一个短暂的通知消息 (Toast)。
 * @param {string} message - 通知消息内容。
 * @param {'success'|'error'} type - 通知类型，默认为 'success'。
 */
function showToast(message, type='success') {
    const toastElement = document.createElement('div');
    toastElement.className = `toast ${type}`;
    toastElement.innerHTML = type === 'error' ?
        `<i class="fa-solid fa-circle-xmark"></i> ${message}` :
        `<i class="fa-solid fa-circle-check"></i> ${message}`;
    $('#toast-container').appendChild(toastElement);
    setTimeout(() => toastElement.remove(), 3000);
}

/**
 * 格式化文件大小。
 * @param {number} bytes - 文件大小（字节）。
 * @returns {string} - 格式化后的文件大小字符串。
 */
const formatSize = (bytes) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

/**
 * 格式化日期字符串为本地时间。
 * @param {string} dateString - 日期字符串。
 * @returns {string} - 格式化后的日期时间字符串。
 */
const formatDate = (dateString) => new Date(dateString).toLocaleString('zh-CN', { hour12: false });
