const state = {
    interviews: [],
    viewInterviews: [],
    filterSuggestion: '',
    pagination: {
        page: 1,
        limit: 10
    },
    editInterviewId: null,
    allResumes: [],
    allCompanyPositions: [],
    allRecruitmentSources: [],
    allInterviewers: [],
    selectedFileIdForForm: null,
    currentFileSearchKeyword: '',
};

const $ = (selector) => document.querySelector(selector);
const $all = (selector) => document.querySelectorAll(selector);

let interviewFormApp; // Vue 实例

document.addEventListener('DOMContentLoaded', () => {
    fetchInterviews();
    fetchResumesForSelect();
    fetchCompanyPositionsForSelect();
    fetchRecruitmentSourcesForSelect();
    fetchInterviewersForSelect();
    setupEvents();
    initVueApp();
    document.addEventListener('keydown', handleEscapeKey);
});

function initVueApp() {
    interviewFormApp = new Vue({
        el: '#interview-form-modal',
        data() {
            return {
                interviewForm: {
                    candidate_name: '',
                    company_position_id: '',
                    phone_number: '',
                    email: '',
                    resume_file_id: '',
                    resume_file_display: '', // 用于显示文件名称
                    recruitment_source_id: '',
                    interview_round: '初试',
                    interview_method: '现场面试',
                    interview_time: '',
                    interviewer_id: '',
                    interview_rating: 0,
                    suggestion: '通过',
                    comments: ''
                },
                // Vue 实例内部保存这些列表，以便在模板中使用 v-for
                allCompanyPositions: state.allCompanyPositions,
                allRecruitmentSources: state.allRecruitmentSources,
                allInterviewers: state.allInterviewers,
            };
        },
        watch: {
            'interviewForm.interview_rating'(newValue) {
                let parsedValue = parseInt(newValue, 10);
                if (isNaN(parsedValue)) {
                    this.interviewForm.interview_rating = 0;
                    return;
                }
                if (parsedValue < 0) {
                    this.interviewForm.interview_rating = 0;
                }
                else if (parsedValue > 10) {
                    this.interviewForm.interview_rating = 10;
                }
            }
        },
        mounted() {
            // 初始化面试时间
            this.interviewForm.interview_time = formatDateForElementUIDisplay(new Date());
            if (this.$refs.interviewTimePicker) {
                this.$refs.interviewTimePicker.$emit('input', this.interviewForm.interview_time);
            }
        },
        methods: {
            openFileSearchModal() {
                window.openFileSearchModal();
            },
            handleFormSubmit() {
                window.handleFormSubmit(); // 调用全局的 handleFormSubmit
            },
            handleLocalResumeUpload(event) {
                window.handleLocalResumeUpload(event); // 调用全局的 handleLocalResumeUpload
            },
            openEditInterviewModalFromDetail(id) {
                window.openEditInterviewModal(id); // 调用全局的 openEditInterviewModal
            }
        }
    });
}

function setupEvents() {
    $('#search-input').addEventListener('input', (e) => {
        state.pagination.page = 1;
        applyFilters(e.target.value);
    });

    $all('#filter-tabs .tab').forEach(btn => {
        btn.addEventListener('click', (e) => {
            $all('.tab').forEach(t => t.classList.remove('active'));
            e.currentTarget.classList.add('active');
            state.filterSuggestion = e.currentTarget.dataset.suggestion;
            state.pagination.page = 1;
            applyFilters($('#search-input').value);
        });
    });

    $('#file-search-input').addEventListener('input', (e) => {
        state.currentFileSearchKeyword = e.target.value;
        renderFileSearchList();
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

window.fetchInterviews = async () => {
    const btnIcon = $('.btn-secondary-pill i');
    if (btnIcon) btnIcon.classList.add('fa-spin');

    try {
        const res = await fetch('/api/v1/interviews?page=1&limit=99999');
        const data = await res.json();
        if (data.success) {
            state.interviews = data.interviews.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
            applyFilters($('#search-input').value);
            updateStats();
        }
    } catch (e) {
        showToast('无法连接服务器或加载面试记录', 'error');
        console.error("Error fetching interviews:", e);
    } finally {
        if (btnIcon) setTimeout(() => btnIcon.classList.remove('fa-spin'), 500);
    }
};

// 添加 handleEscapeKey 函数
function handleEscapeKey(event) {
    if (event.key === 'Escape') {
        // 按照模态框的层级/优先级来检查和关闭
        if ($('#confirm-modal').style.display === 'flex') {
            // 如果确认模态框打开，Esc键应视为取消
            $('#btn-cancel-confirm').click();
            event.preventDefault(); // 阻止Esc键的默认行为
            return;
        }
        if ($('#file-search-modal').style.display === 'flex') {
            closeModal('file-search-modal');
            event.preventDefault();
            return;
        }
        if ($('#interview-form-modal').style.display === 'flex') {
            closeModal('interview-form-modal');
            event.preventDefault();
            return;
        }
        if ($('#detail-modal').style.display === 'flex') {
            closeModal('detail-modal');
            event.preventDefault();
            return;
        }
    }
}

window.fetchResumesForSelect = async () => {
    try {
        const res = await fetch('/api/v1/files?page=1&limit=99999');
        const data = await res.json();
        if (data.success) {
            state.allResumes = data.files.filter(f => ['doc', 'docx', 'pdf'].includes(f.extension.toLowerCase()));
        }
    } catch (e) {
        showToast('无法加载简历文件列表', 'error');
        console.error("Error fetching resumes:", e);
    }
};

window.fetchCompanyPositionsForSelect = async () => {
    try {
        const res = await fetch('/api/v1/company-positions?page=1&limit=999');
        const data = await res.json();
        if (data.success) {
            state.allCompanyPositions = data.positions;
            if (interviewFormApp) interviewFormApp.allCompanyPositions = data.positions;
        }
    } catch (e) {
        showToast('无法加载应聘岗位列表', 'error');
        console.error("Error fetching company positions:", e);
    }
};

window.fetchRecruitmentSourcesForSelect = async () => {
    try {
        const res = await fetch('/api/v1/recruitment-sources?page=1&limit=999');
        const data = await res.json();
        if (data.success) {
            state.allRecruitmentSources = data.sources;
            if (interviewFormApp) interviewFormApp.allRecruitmentSources = data.sources;
        }
    } catch (e) {
        showToast('无法加载招聘来源列表', 'error');
        console.error("Error fetching recruitment sources:", e);
    }
};

window.fetchInterviewersForSelect = async () => {
    try {
        const res = await fetch('/api/v1/employees?page=1&limit=999');
        const data = await res.json();
        if (data.success) {
            state.allInterviewers = data.employees;
            if (interviewFormApp) interviewFormApp.allInterviewers = data.employees;
        }
    } catch (e) {
        showToast('无法加载面试官列表', 'error');
        console.error("Error fetching interviewers:", e);
    }
};

/**
 * 根据文件ID更新面试表单中的简历文件显示。
 * @param {string} fileId - 选中的简历文件ID。
 */
function updateResumeDisplay(fileId) {
    const selectedFile = state.allResumes.find(f => f.id === fileId);
    if (selectedFile) {
        interviewFormApp.interviewForm.resume_file_display = `${selectedFile.original_name} (${formatDate(selectedFile.created_at)})`;
        interviewFormApp.interviewForm.resume_file_id = selectedFile.id;
    } else {
        interviewFormApp.interviewForm.resume_file_display = '';
        interviewFormApp.interviewForm.resume_file_id = '';
    }
}

/**
 * 处理本地简历文件上传。
 * @param {Event} event - 文件输入框的 change 事件。
 */
window.handleLocalResumeUpload = async (event) => {
    const files = event.target.files;
    if (files.length === 0) return;

    const file = files[0];
    const formData = new FormData();
    formData.append('file', file);

    try {
        showToast('简历文件上传中...', 'info');
        const res = await fetch('/api/v1/interviews/upload-resume', {
            method: 'POST',
            body: formData
        });
        const data = await res.json();
        if (data.success && data.file) {
            showToast('简历文件上传成功');
            state.allResumes.push(data.file); // 将新上传的文件添加到简历列表中
            updateResumeDisplay(data.file.id); // 更新表单显示
        } else {
            showToast('简历文件上传失败: ' + (data.message || '未知错误'), 'error');
        }
    } catch (e) {
        showToast('网络错误或服务器无响应', 'error');
        console.error("Local resume upload error:", e);
    } finally {
        event.target.value = ''; // 清空文件输入框，以便再次选择相同文件
    }
};

/**
 * 根据关键字和建议类型过滤面试记录。
 * @param {string} keyword - 搜索关键字。
 */
function applyFilters(keyword = '') {
    const key = keyword.toLowerCase();
    state.viewInterviews = state.interviews.filter(item => {
        const companyPositionName = item.company_position ? item.company_position.name.toLowerCase() : '';
        const recruitmentSourceName = item.recruitment_source ? item.recruitment_source.name.toLowerCase() : '';
        const interviewerName = item.interviewer ? item.interviewer.name.toLowerCase() : '';

        const keywordMatch = item.candidate_name.toLowerCase().includes(key) ||
            companyPositionName.includes(key) ||
            interviewerName.includes(key) ||
            item.phone_number.toLowerCase().includes(key) ||
            item.email.toLowerCase().includes(key) ||
            recruitmentSourceName.includes(key);
        const suggestionMatch = !state.filterSuggestion || item.suggestion === state.filterSuggestion;
        return keywordMatch && suggestionMatch;
    });
    renderTable();
    renderPagination();
}

/**
 * 获取当前页的经过过滤的面试记录数据。
 * @returns {Array} - 当前页的面试记录数据。
 */
function getPaginatedData() {
    const start = (state.pagination.page - 1) * state.pagination.limit;
    return state.viewInterviews.slice(start, start + state.pagination.limit);
}

/**
 * 渲染面试记录表格。
 */
function renderTable() {
    const list = getPaginatedData();
    const tbody = $('#interview-list-body');
    tbody.innerHTML = '';

    if (!list.length) {
        tbody.innerHTML = `<tr><td colspan="10" class="text-center" style="padding:40px;color:#9ca3af;">没有找到匹配的面试记录</td></tr>`;
        return;
    }

    list.forEach((item, index) => {
        const tr = document.createElement('tr');
        tr.innerHTML = `
            <td>${(state.pagination.page - 1) * state.pagination.limit + index + 1}</td>
            <td>
                <span class="main-text">${item.candidate_name}</span>
                <br>
                <span class="sub-text">${item.company_position ? item.company_position.name : 'N/A'}</span>
            </td>
            <td>
                <span class="main-text">${item.phone_number}</span>
                <br>
                <span class="sub-text">${item.email}</span>
            </td>
            <td>${item.recruitment_source ? item.recruitment_source.name : 'N/A'}</td>
            <td><span class="badge badge-${getRoundBadgeClass(item.interview_round)}">${item.interview_round}</span></td>
            <td>${formatDateTime(item.interview_time)}</td>
            <td>${item.interviewer ? item.interviewer.name : 'N/A'}</td>
            <td class="text-center">${item.interview_rating || '0'}</td>
            <td class="text-center"><span class="badge badge-${getSuggestionBadgeClass(item.suggestion)}">${item.suggestion}</span></td>
            <td class="text-right">
                <button class="action-btn btn-view" onclick="openDetailModal('${item.id}')" title="详情"><i class="fa-solid fa-eye"></i></button>
                <button class="action-btn btn-edit" onclick="openEditInterviewModal('${item.id}')" title="编辑"><i class="fa-solid fa-edit"></i></button>
                <button class="action-btn btn-del" onclick="deleteInterview('${item.id}')" title="删除"><i class="fa-solid fa-trash-can"></i></button>
            </td>
        `;
        tbody.appendChild(tr);
    });
}

/**
 * 渲染分页控件。
 */
function renderPagination() {
    const totalItemsCount = state.viewInterviews.length;
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
    if (isEnabled) button.onclick = onClickCallback;
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
 * 打开新增面试记录弹窗并初始化表单。
 */
window.openCreateInterviewModal = () => {
    state.editInterviewId = null;
    $('#form-modal-title').innerText = '新增面试记录';
    $('#interview-form-modal .modal-tag').innerText = 'New';
    $('#interview-form-modal .modal-tag').className = 'modal-tag'; // 清除编辑状态的tag颜色

    Object.assign(interviewFormApp.interviewForm, {
        candidate_name: '',
        company_position_id: '',
        phone_number: '',
        email: '',
        resume_file_id: '',
        resume_file_display: '',
        recruitment_source_id: '',
        interview_round: '初试',
        interview_method: '现场面试',
        interview_time: formatDateForElementUIDisplay(new Date()),
        interviewer_id: '',
        interview_rating: 0,
        suggestion: '通过',
        comments: ''
    });

    if (interviewFormApp.$refs.interviewTimePicker) {
        interviewFormApp.$refs.interviewTimePicker.$emit('input', interviewFormApp.interviewForm.interview_time);
    }

    $('#interview-form-modal').style.display = 'flex';
};

/**
 * 打开编辑面试记录弹窗并填充表单数据。
 * @param {string} id - 要编辑的面试记录ID。
 */
window.openEditInterviewModal = (id) => {
    const interview = state.interviews.find(item => item.id === id);
    if (!interview) {
        showToast('面试记录未找到', 'error');
        return;
    }

    state.editInterviewId = id;
    $('#form-modal-title').innerText = '编辑面试记录';
    $('#interview-form-modal .modal-tag').innerText = 'Edit';
    $('#interview-form-modal .modal-tag').className = 'modal-tag primary'; // 设置编辑状态的tag颜色

    Object.assign(interviewFormApp.interviewForm, {
        candidate_name: interview.candidate_name,
        company_position_id: interview.company_position_id || '', // 使用或空字符串确保 v-model 不出错
        phone_number: interview.phone_number,
        email: interview.email,
        recruitment_source_id: interview.recruitment_source_id || '',
        interview_round: interview.interview_round,
        interview_method: interview.interview_method,
        interview_time: formatDateForElementUIDisplay(new Date(interview.interview_time)),
        interviewer_id: interview.interviewer_id || '',
        interview_rating: interview.interview_rating,
        suggestion: interview.suggestion,
        comments: interview.comments
    });

    const selectedFile = state.allResumes.find(f => f.id === interview.resume_file_id);
    interviewFormApp.interviewForm.resume_file_id = interview.resume_file_id || '';
    interviewFormApp.interviewForm.resume_file_display = selectedFile ? `${selectedFile.original_name} (${formatDate(selectedFile.created_at)})` : '';

    if (interviewFormApp.$refs.interviewTimePicker) {
        interviewFormApp.$refs.interviewTimePicker.$emit('input', interviewFormApp.interviewForm.interview_time);
    }

    $('#interview-form-modal').style.display = 'flex';
};

/**
 * 处理面试表单的提交（新增或更新）。
 */
window.handleFormSubmit = async () => {
    const data = { ...interviewFormApp.interviewForm };
    const isEditMode = !!state.editInterviewId;

    if (!data.resume_file_id) {
        showToast('请选择简历文件', 'error');
        return;
    }

    // 如果是新增模式，检查简历是否已被使用
    if (!isEditMode) {
        const isResumeUsed = state.interviews.some(interview =>
            interview.resume_file_id === data.resume_file_id
        );
        if (isResumeUsed) {
            showToast('该简历文件已被其他面试记录使用，请选择其他简历文件', 'error');
            return;
        }
    }

    if (!data.company_position_id) {
        showToast('请选择应聘岗位', 'error');
        return;
    }
    if (!data.recruitment_source_id) {
        showToast('请选择招聘来源', 'error');
        return;
    }
    if (!data.interviewer_id) {
        showToast('请选择面试官', 'error');
        return;
    }

    if (data.interview_time) {
        let localDateTimeStr = data.interview_time.replace(' ', 'T');
        try {
            const dateObj = new Date(localDateTimeStr);
            if (isNaN(dateObj.getTime())) {
                throw new Error("Invalid date string from Element UI");
            }
            // 转换为 UTC ISO 格式，后端通常期望这种格式
            data.interview_time = dateObj.toISOString().slice(0, 19) + 'Z';
        } catch (e) {
            console.error("Error converting interview_time to ISO format:", e);
            showToast('面试时间格式错误，请重新选择', 'error');
            return;
        }
    } else {
        showToast('请选择面试时间', 'error');
        return;
    }

    const method = state.editInterviewId ? 'PUT' : 'POST';
    const url = state.editInterviewId ? `/api/v1/interviews/${state.editInterviewId}` : '/api/v1/interviews';

    try {
        const res = await fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });
        const result = await res.json();
        if (result.success) {
            showToast(state.editInterviewId ? '面试记录更新成功' : '面试记录创建成功');
            closeModal('interview-form-modal');
            fetchInterviews();
            fetchResumesForSelect(); // 重新加载简历以确保新增的本地文件可见
            fetchCompanyPositionsForSelect();
            fetchRecruitmentSourcesForSelect();
            fetchInterviewersForSelect();
        } else {
            showToast(`操作失败: ${result.message}`, 'error');
        }
    } catch (err) {
        showToast('网络错误或服务器无响应', 'error');
        console.error("Form submit error:", err);
    }
};

/**
 * 删除面试记录。
 * @param {string} id - 要删除的面试记录ID。
 */
window.deleteInterview = async (id) => {
    // 明确提示用户这是永久删除
    const confirmed = await showConfirmModal('确认删除?', '此面试记录将被<b>永久删除</b>，且无法恢复。您确定要继续吗？');
    if (!confirmed) return;
    try {
        const res = await fetch(`/api/v1/interviews/${id}`, { method: 'DELETE' });
        if (res.ok) {
            showToast('面试记录已成功删除');
            fetchInterviews();
        } else {
            const errorData = await res.json();
            showToast(`删除失败：${errorData.message || '未知错误'}`, 'error');
        }
    } catch (e) {
        showToast('删除操作发生网络错误', 'error');
        console.error("Delete interview error:", e);
    }
};

/**
 * 打开面试详情弹窗并显示记录数据。
 * @param {string} id - 要查看的面试记录ID。
 */
window.openDetailModal = (id) => {
    const interview = state.interviews.find(item => item.id === id);
    if (!interview) {
        showToast('面试记录未找到', 'error');
        return;
    }

    state.editInterviewId = id; // 记录当前详情页的ID，以便编辑按钮使用
    $('#detail-modal-round-tag').innerText = interview.interview_round;
    $('#detail-modal-round-tag').className = `modal-tag ${getRoundBadgeClass(interview.interview_round)}`;

    const resumeFileName = interview.resume_file ? interview.resume_file.original_name : '未关联简历';
    const resumeFileLink = interview.resume_file ? `/api/v1/files/${interview.resume_file.id}/download` : '#';
    const resumeFileDisplay = interview.resume_file ? `<a href="${resumeFileLink}" target="_blank" class="detail-link">${resumeFileName} <i class="fa-solid fa-download"></i></a>` : resumeFileName;

    const detailContentHtml = `
        <div class="detail-key">姓名</div><div class="detail-val">${interview.candidate_name}</div>
        <div class="detail-key">手机号</div><div class="detail-val">${interview.phone_number}</div>
        <div class="detail-key">邮箱</div><div class="detail-val">${interview.email}</div>
        <div class="detail-key">应聘岗位</div><div class="detail-val">${interview.company_position ? interview.company_position.name : 'N/A'}</div>
        <div class="detail-key">简历文件</div><div class="detail-val">${resumeFileDisplay}</div>
        <div class="detail-key">招聘来源</div><div class="detail-val">${interview.recruitment_source ? interview.recruitment_source.name : 'N/A'}</div>
        <div class="detail-key">面试轮次</div><div class="detail-val">${interview.interview_round}</div>
        <div class="detail-key">面试形式</div><div class="detail-val">${interview.interview_method}</div>
        <div class="detail-key">面试时间</div><div class="detail-val">${formatDateTime(interview.interview_time)}</div>
        <div class="detail-key">面试官</div><div class="detail-val">${interview.interviewer ? interview.interviewer.name : 'N/A'}</div>
        <div class="detail-key">面试评分</div><div class="detail-val">${interview.interview_rating || '0'}</div>
        <div class="detail-key">面试评价</div><div class="detail-val">${interview.comments || '无'}</div>
        <div class="detail-key">本轮建议</div><div class="detail-val">${interview.suggestion}</div>
        <div class="detail-key">创建时间</div><div class="detail-val">${formatDateTime(interview.created_at)}</div>
    `;
    $('#detail-content').innerHTML = detailContentHtml;
    $('#detail-modal').style.display = 'flex';
};

/**
 * 关闭指定ID的模态框。
 * @param {string} id - 模态框的ID。
 */
window.closeModal = (id) => {
    $('#' + id).style.display = 'none';
    if (id === 'interview-form-modal') {
        // 重置 select 选项，确保下次打开是最新数据
        fetchCompanyPositionsForSelect();
        fetchRecruitmentSourcesForSelect();
        fetchInterviewersForSelect();
    }
};

/**
 * 打开简历文件搜索弹窗并初始化状态。
 */
window.openFileSearchModal = () => {
    state.selectedFileIdForForm = interviewFormApp.interviewForm.resume_file_id;
    state.currentFileSearchKeyword = '';
    $('#file-search-input').value = '';
    renderFileSearchList();
    $('#file-search-modal').style.display = 'flex';
};

/**
 * 渲染简历文件搜索列表。
 */
function renderFileSearchList() {
    const listContainer = $('#file-search-list');
    listContainer.innerHTML = '';
    const key = state.currentFileSearchKeyword.toLowerCase();

    // 获取所有已使用的简历文件ID
    const usedResumeIds = new Set();
    state.interviews.forEach(interview => {
        if (interview.resume_file_id) {
            usedResumeIds.add(interview.resume_file_id);
        }
    });

    const filteredFiles = state.allResumes.filter(f => {
        return f.original_name.toLowerCase().includes(key) || f.extension.toLowerCase().includes(key);
    });

    if (filteredFiles.length === 0) {
        listContainer.innerHTML = `<div class="text-center" style="padding:20px;color:#9ca3af;">没有找到匹配的简历文件</div>`;
        return;
    }

    filteredFiles.forEach(file => {
        const isUsed = usedResumeIds.has(file.id);
        const isSelected = file.id === state.selectedFileIdForForm;

        let itemClassName = 'file-search-item';
        if (isSelected) itemClassName += ' active';
        if (isUsed) itemClassName += ' used-file';

        const item = document.createElement('div');
        item.className = itemClassName;
        item.dataset.id = file.id;
        item.title = isUsed ? '该简历已被其他面试记录使用' : '';

        const ext = file.extension.toLowerCase();
        const iconInfo = ext === 'pdf' ? { iconClass: 'fa-file-pdf', badgeClass: 'fb-pdf' } : { iconClass: 'fa-file-word', badgeClass: 'fb-word' };

        item.innerHTML = `
            <div class="file-badge ${iconInfo.badgeClass}"><i class="fa-solid ${iconInfo.iconClass}"></i></div>
            <div class="file-search-item-info">
                <span class="file-search-item-name" style="${isUsed ? 'color:#999;' : ''}">
                    ${file.original_name}
                    ${isUsed ? '<span style="color:#ff6b6b;font-size:10px;margin-left:5px;">(已使用)</span>' : ''}
                </span>
                <span class="file-search-item-meta">${formatSize(file.size)} | ${formatDate(file.created_at)}</span>
            </div>
        `;

        // 如果是已使用的简历，点击时有不同提示
        item.addEventListener('click', () => {
            if (isUsed) {
                showToast('该简历已被其他面试记录使用，请选择其他简历', 'warning');
                return;
            }
            $all('.file-search-item').forEach(el => el.classList.remove('active'));
            item.classList.add('active');
            state.selectedFileIdForForm = file.id;
        });

        item.addEventListener('dblclick', () => {
            if (isUsed) {
                showToast('该简历已被其他面试记录使用，请选择其他简历', 'warning');
                return;
            }
            state.selectedFileIdForForm = file.id;
            selectResumeFile();
        });

        listContainer.appendChild(item);
    });
}

/**
 * 确认并选择简历文件。
 */
window.selectResumeFile = () => {
    if (state.selectedFileIdForForm) {
        updateResumeDisplay(state.selectedFileIdForForm);
        closeModal('file-search-modal');
    } else {
        showToast('请选择一个简历文件', 'error');
    }
};

/**
 * 更新页面上的面试统计数据。
 */
function updateStats() {
    $('#stat-total-interviews').innerText = state.interviews.length;
    const pendingCount = state.interviews.filter(i => i.suggestion === '待定').length;
    $('#stat-pending-interviews').innerText = pendingCount;
    const passedCount = state.interviews.filter(i => i.suggestion === '通过').length;
    $('#stat-passed-interviews').innerText = passedCount;
}

/**
 * 显示一个短暂的通知消息 (Toast)。
 * @param {string} message - 通知消息内容。
 * @param {'success'|'error'|'info'} type - 通知类型，默认为 'success'。
 */
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
 * 格式化日期字符串为本地日期。
 * @param {string} dateString - 日期字符串。
 * @returns {string} - 格式化后的日期字符串。
 */
const formatDate = (dateString) => {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) {
        return '无效日期';
    }
    return date.toLocaleString('zh-CN', { hour12: false, year: 'numeric', month: '2-digit', day: '2-digit' });
};

/**
 * 格式化日期字符串为本地日期时间。
 * @param {string} dateString - 日期字符串。
 * @returns {string} - 格式化后的日期时间字符串。
 */
const formatDateTime = (dateString) => {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) {
        return '无效日期';
    }
    return date.toLocaleString('zh-CN', { hour12: false });
};

/**
 * 格式化日期对象为 Element UI date-picker 预期的字符串格式 (yyyy-MM-dd HH:mm:ss)。
 * @param {Date} date - 日期对象。
 * @returns {string} - 格式化后的日期时间字符串。
 */
const formatDateForElementUIDisplay = (date) => {
    const year = date.getFullYear();
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const day = date.getDate().toString().padStart(2, '0');
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    const seconds = date.getSeconds().toString().padStart(2, '0');
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
};

/**
 * 根据面试轮次返回对应的徽章颜色类。
 * @param {string} round - 面试轮次 (e.g., '初试', '复试', '终面')。
 * @returns {string} - 徽章颜色类名。
 */
function getRoundBadgeClass(round) {
    switch (round) {
        case '初试': return 'blue';
        case '复试': return 'yellow';
        case '终面': return 'green';
        default: return 'gray';
    }
}

/**
 * 根据面试建议返回对应的徽章颜色类。
 * @param {string} suggestion - 面试建议 (e.g., '通过', '待定', '淘汰')。
 * @returns {string} - 徽章颜色类名。
 */
function getSuggestionBadgeClass(suggestion) {
    switch (suggestion) {
        case '通过': return 'green';
        case '待定': return 'yellow';
        case '淘汰': return 'red';
        default: return 'gray';
    }
}

/**
 * 关闭详情弹窗并打开编辑表单
 * @param {string} id - 面试记录ID
 */
window.closeDetailAndEdit = (id) => {
    // 先关闭详情弹窗
    closeModal('detail-modal');
    // 然后打开编辑表单
    setTimeout(() => {
        openEditInterviewModal(id);
    }, 100); // 小延迟确保详情弹窗已关闭
};
