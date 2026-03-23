const state = {
    items: [],
    viewItems: [],
    pagination: {
        page: 1,
        limit: 10
    },
    editItemId: null,
    currentSearchKeyword: '',
};

const $ = (selector) => document.querySelector(selector);
const $all = (selector) => document.querySelectorAll(selector);

let itemFormApp; // Vue 实例

document.addEventListener('DOMContentLoaded', () => {
    fetchItems();
    setupEvents();
    initVueApp();
    document.addEventListener('keydown', handleEscapeKey);
});

function initVueApp() {
    itemFormApp = new Vue({
        el: '#form-modal',
        data() {
            return {
                itemForm: {
                    name: ''
                }
            };
        },
        methods: {
            handleFormSubmit() {
                window.handleFormSubmit();
            }
        }
    });
}

function setupEvents() {
    $('#search-input').addEventListener('input', (e) => {
        state.pagination.page = 1;
        state.currentSearchKeyword = e.target.value;
        applyFilters();
    });
}

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

window.fetchItems = async () => {
    const btnIcon = $('.btn-secondary-pill i');
    if (btnIcon) btnIcon.classList.add('fa-spin');

    try {
        const res = await fetch('/api/v1/recruitment-sources?page=1&limit=99999');
        const data = await res.json();
        if (data.success) {
            state.items = data.sources.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
            applyFilters();
            updateStats();
        }
    } catch (e) {
        showToast('无法连接服务器或加载招聘来源列表', 'error');
        console.error("Error fetching items:", e);
    } finally {
        if (btnIcon) setTimeout(() => btnIcon.classList.remove('fa-spin'), 500);
    }
};

// 添加 handleEscapeKey 函数
function handleEscapeKey(event) {
    if (event.key === 'Escape') {
        if ($('#confirm-modal').style.display === 'flex') {
            $('#btn-cancel-confirm').click();
            event.preventDefault();
            return;
        }
        if ($('#form-modal').style.display === 'flex') {
            closeModal('form-modal');
            event.preventDefault();
            return;
        }
    }
}

function applyFilters() {
    const key = state.currentSearchKeyword.toLowerCase();
    state.viewItems = state.items.filter(item => {
        return item.name.toLowerCase().includes(key);
    });
    renderTable();
    renderPagination();
}

function getPaginatedData() {
    const start = (state.pagination.page - 1) * state.pagination.limit;
    return state.viewItems.slice(start, start + state.pagination.limit);
}

function renderTable() {
    const list = getPaginatedData();
    const tbody = $('#list-body');
    tbody.innerHTML = '';

    if (!list.length) {
        tbody.innerHTML = `<tr><td colspan="4" class="text-center" style="padding:40px;color:#9ca3af;">没有找到匹配的招聘来源</td></tr>`;
        return;
    }

    list.forEach((item, index) => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td>${(state.pagination.page - 1) * state.pagination.limit + index + 1}</td>
            <td><span class="main-text">${item.name}</span></td>
            <td>${formatDateTime(item.created_at)}</td>
            <td>
                <button class="action-btn btn-edit" onclick="openEditItemModal('${item.id}')" title="编辑"><i class="fa-solid fa-edit"></i></button>
                <button class="action-btn btn-del" onclick="deleteItem('${item.id}')" title="删除"><i class="fa-solid fa-trash-can"></i></button>
            </td>
        `;
        tbody.appendChild(tr);
    });
}

function renderPagination() {
    const totalItemsCount = state.viewItems.length;
    const itemsPerPage = state.pagination.limit;
    const totalPages = Math.ceil(totalItemsCount / itemsPerPage);
    const currentPage = state.pagination.page;
    const paginationContainer = $('#pagination');

    paginationContainer.innerHTML = '';
    if (totalItemsCount === 0 || (totalPages <= 1 && totalItemsCount <= itemsPerPage)) {
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

function createPageBtn(htmlContent, isEnabled, onClickCallback) {
    const button = document.createElement('button');
    button.className = 'page-btn';
    button.innerHTML = htmlContent;
    button.disabled = !isEnabled;
    if (isEnabled) button.onclick = onClickCallback;
    return button;
}

function changePage(newPage) {
    state.pagination.page = newPage;
    renderTable();
    renderPagination();
}

window.openCreateItemModal = () => {
    state.editItemId = null;
    $('#form-modal-title').innerText = '新增招聘来源';
    $('#form-modal .modal-tag').innerText = 'New';
    $('#form-modal .modal-tag').className = 'modal-tag';

    Object.assign(itemFormApp.itemForm, {
        name: ''
    });

    $('#form-modal').style.display = 'flex';
};

window.openEditItemModal = (id) => {
    const item = state.items.find(i => i.id == id);
    if (!item) {
        showToast('招聘来源未找到', 'error');
        return;
    }

    state.editItemId = id;
    $('#form-modal-title').innerText = '编辑招聘来源';
    $('#form-modal .modal-tag').innerText = 'Edit';
    $('#form-modal .modal-tag').className = 'modal-tag primary';

    Object.assign(itemFormApp.itemForm, {
        name: item.name
    });

    $('#form-modal').style.display = 'flex';
};

window.handleFormSubmit = async () => {
    const data = { ...itemFormApp.itemForm };

    if (!data.name) {
        showToast('招聘来源名称不能为空', 'error');
        return;
    }

    const method = state.editItemId ? 'PUT' : 'POST';
    const url = state.editItemId ? `/api/v1/recruitment-sources/${state.editItemId}` : '/api/v1/recruitment-sources';

    try {
        const res = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        const result = await res.json();
        if (result.success) {
            showToast(state.editItemId ? '招聘来源更新成功' : '招聘来源创建成功');
            closeModal('form-modal');
            fetchItems();
        } else {
            showToast(`操作失败: ${result.message}`, 'error');
        }
    } catch (err) {
        showToast('网络错误或服务器无响应', 'error');
        console.error("Form submit error:", err);
    }
};

window.deleteItem = async (id) => {
    const confirmed = await showConfirmModal('确认删除?', '此招聘来源将被<b>永久删除</b>，且无法恢复。您确定要继续吗？');
    if (!confirmed) return;
    try {
        const res = await fetch(`/api/v1/recruitment-sources/${id}`, { method: 'DELETE' });
        if (res.ok) {
            showToast('招聘来源已成功删除');
            fetchItems();
        } else {
            const errorData = await res.json();
            showToast(`删除失败：${errorData.message || '未知错误'}`, 'error');
        }
    } catch (e) {
        showToast('删除操作发生网络错误', 'error');
        console.error("Delete item error:", e);
    }
};

window.closeModal = (id) => {
    $('#' + id).style.display = 'none';
};

function updateStats() {
    const today = formatDate(new Date());
    const todayNewCount = state.items.filter(item => formatDate(item.created_at) === today).length;
    $('#stat-total-items').innerText = state.items.length;
    $('#stat-active-items').innerText = state.items.length;
    $('#stat-today-new').innerText = todayNewCount;
}

function showToast(message, type = 'success') {
    const toastElement = document.createElement('div');
    toastElement.className = `toast ${type}`;
    toastElement.innerHTML = type === 'error' ?
        `<i class="fa-solid fa-circle-xmark"></i> ${message}` :
        type === 'info' ?
            `<i class="fa-solid fa-circle-info"></i> ${message}` :
            `<i class="fa-solid fa-circle-check"></i> ${message}`;
    $('#toast-container').appendChild(toastElement);
    setTimeout(() => toastElement.remove(), 3000);
}

const formatDate = (dateString) => {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) {
        return '无效日期';
    }
    return date.toLocaleString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' });
};

const formatDateTime = (dateString) => {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) {
        return '无效日期';
    }
    return date.toLocaleString('zh-CN', { hour12: false });
};
