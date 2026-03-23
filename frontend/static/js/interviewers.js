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
                    name: '',
                    email: '',
                    department: ''
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
        const res = await fetch('/api/v1/employees?page=1&limit=99999');
        const data = await res.json();
        if (data.success) {
            state.items = data.employees.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
            applyFilters();
            updateStats();
        }
    } catch (e) {
        showToast('无法连接服务器或加载面试官列表', 'error');
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
        return item.name.toLowerCase().includes(key) ||
            item.email.toLowerCase().includes(key) ||
            (item.department && item.department.toLowerCase().includes(key));
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
        tbody.innerHTML = `<tr><td colspan="6" class="text-center" style="padding:40px;color:#9ca3af;">没有找到匹配的面试官</td></tr>`;
        return;
    }

    list.forEach((item, index) => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td>${(state.pagination.page - 1) * state.pagination.limit + index + 1}</td>
            <td><span class="main-text">${item.name}</span></td>
            <td><span class="sub-text">${item.email}</span></td>
            <td><span class="sub-text">${item.department || '无'}</span></td>
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
    $('#form-modal-title').innerText = '新增面试官';
    $('#form-modal .modal-tag').innerText = 'New';
    $('#form-modal .modal-tag').className = 'modal-tag';

    Object.assign(itemFormApp.itemForm, {
        name: '',
        email: '',
        department: ''
    });

    $('#form-modal').style.display = 'flex';
};

window.openEditItemModal = (id) => {
    const item = state.items.find(i => i.id == id);
    if (!item) {
        showToast('面试官未找到', 'error');
        return;
    }

    state.editItemId = id;
    $('#form-modal-title').innerText = '编辑面试官';
    $('#form-modal .modal-tag').innerText = 'Edit';
    $('#form-modal .modal-tag').className = 'modal-tag primary';

    Object.assign(itemFormApp.itemForm, {
        name: item.name,
        email: item.email,
        department: item.department || '' // Ensure department is not null if unset
    });

    $('#form-modal').style.display = 'flex';
};

window.handleFormSubmit = async () => {
    const data = { ...itemFormApp.itemForm };

    if (!data.name) {
        showToast('姓名不能为空', 'error');
        return;
    }
    if (!data.email || !isValidEmail(data.email)) {
        showToast('请输入有效的邮箱', 'error');
        return;
    }

    const method = state.editItemId ? 'PUT' : 'POST';
    const url = state.editItemId ? `/api/v1/employees/${state.editItemId}` : '/api/v1/employees';

    try {
        const res = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        const result = await res.json();
        if (result.success) {
            showToast(state.editItemId ? '面试官信息更新成功' : '面试官创建成功');
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
    const confirmed = await showConfirmModal('确认删除?', '此面试官信息将被<b>永久删除</b>，且无法恢复。您确定要继续吗？');
    if (!confirmed) return;
    try {
        const res = await fetch(`/api/v1/employees/${id}`, { method: 'DELETE' });
        if (res.ok) {
            showToast('面试官已成功删除');
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
    const total = state.items.length;
    const validEmails = state.items.filter(item => isValidEmail(item.email)).length;
    const withDepartment = state.items.filter(item => item.department && item.department.trim() !== '').length;

    $('#stat-total-items').innerText = total;
    $('#stat-valid-emails').innerText = validEmails;
    $('#stat-with-department').innerText = withDepartment;
}

function isValidEmail(email) {
    // 简单的邮箱验证
    return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
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
