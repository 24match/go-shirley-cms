let token = localStorage.getItem('token');
let bannerUploadedFile = null, aboutUploadedFile = null, productUploadedFiles = [], contentUploadedFile = null;

function setToken(newToken) { token = newToken; newToken ? localStorage.setItem('token', newToken) : localStorage.removeItem('token'); }
function getToken() { return token; }

async function apiFetch(url, options = {}) {
  if (!token) { handleAuthError(); throw new Error('No token'); }
  const res = await fetch(url, { ...options, headers: { 'Authorization': 'Bearer ' + token, ...options.headers } });
  if (res.status === 401) { handleAuthError(); throw new Error('Unauthorized'); }
  return res;
}

function handleAuthError() { setToken(null); showToast('登录已过期，请重新登录', 'error'); setTimeout(() => { document.getElementById('adminPage').style.display = 'none'; document.getElementById('loginPage').style.display = 'flex'; }, 1500); }
function showToast(message, type = 'success') { const toast = document.getElementById('toast'); toast.textContent = message; toast.className = 'toast ' + type; toast.classList.add('show'); setTimeout(() => toast.classList.remove('show'), 3000); }

async function login() {
  const username = document.getElementById('username').value;
  const password = document.getElementById('password').value;
  const res = await fetch('/api/login', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ username, password }) });
  const responseData = await res.json();
  if (res.ok) {
    const newToken = responseData.data?.token;
    if (newToken) { setToken(newToken); document.getElementById('loginPage').style.display = 'none'; document.getElementById('adminPage').style.display = 'flex'; loadAllConfigs(); }
    else showToast('登录失败: 未获取到 token', 'error');
  } else showToast('登录失败: ' + (responseData.message || '未知错误'), 'error');
}

function logout() { setToken(null); document.getElementById('adminPage').style.display = 'none'; document.getElementById('loginPage').style.display = 'flex'; }

async function loadAllConfigs() {
  try {
    const [modulesRes, configRes] = await Promise.all([apiFetch('/api/admin/modules'), apiFetch('/api/admin/config')]);
    const modules = (await modulesRes.json()).data || [];
    const configs = (await configRes.json()).data || [];
    const moduleMap = {}; modules.forEach(m => moduleMap[m.moduleName] = m);
    const configMap = {}; configs.forEach(c => { try { configMap[c.pageName] = JSON.parse(c.configData); } catch(e) { configMap[c.pageName] = {}; } });
    
    loadBannerConfig(moduleMap['banner'] || configMap['banner'] || {});
    loadAboutConfig(moduleMap['about'] || configMap['about'] || {});
    loadProductsConfig(moduleMap['products'] || configMap['products'] || {});
    loadFactoryConfig(moduleMap['factory'] || {});
    loadAdvantageConfig(moduleMap['advantage'] || {});
    loadEventsConfig(moduleMap['events'] || configMap['event'] || {});
    loadContactConfig(moduleMap['contact'] || configMap['contact'] || {});
    loadSystemConfig(configMap['brand'], configMap['stats']);
    loadSiteSettings();
  } catch(e) { console.error('Error loading configs:', e); }
}

async function saveModuleConfig(moduleName, formData) {
  const res = await apiFetch('/api/admin/modules', { method: 'POST', body: formData });
  if (res.ok) { showToast('配置保存成功！'); return true; }
  const err = await res.json(); showToast('保存失败: ' + (err.error || '未知错误'), 'error'); return false;
}

async function savePageConfig(pageName, configData) {
  const res = await apiFetch('/api/admin/config', { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ pageName, configData: JSON.stringify(configData) }) });
  if (res.ok) { showToast('配置保存成功！'); return true; }
  showToast('保存失败', 'error'); return false;
}

function handleBannerImageUpload(input) {
  const file = input.files[0]; if (!file) return;
  const previewContainer = document.getElementById('bannerImagesPreview');
  const emptySlot = previewContainer.querySelector('.empty');
  if (!emptySlot) { showToast('Banner模块仅允许上传1张图片', 'error'); return; }
  bannerUploadedFile = file;
  const reader = new FileReader();
  reader.onload = e => {
    emptySlot.innerHTML = `<img src="${e.target.result}" alt="Banner"><div class="image-actions"><button class="action-btn replace" onclick="replaceBannerImage(event)">🔄</button><button class="action-btn delete" onclick="removeBannerImage()">🗑️</button></div>`;
    emptySlot.classList.remove('empty');
    document.getElementById('bannerImageCount').textContent = '1';
    showToast('图片上传成功');
  };
  reader.readAsDataURL(file);
}

function replaceBannerImage(event) { event.stopPropagation(); document.getElementById('bannerImageInput').click(); }

async function removeBannerImage() {
  if (!confirm('确定要删除这张Banner图片吗？')) return;
  const previewContainer = document.getElementById('bannerImagesPreview');
  const imageItem = previewContainer.querySelector('.image-preview-item:not(.empty)');
  if (imageItem) {
    const res = await apiFetch('/api/admin/modules/banner/image', { method: 'DELETE' });
    if (!res.ok) { showToast('删除图片失败', 'error'); return; }
    previewContainer.innerHTML = `<div class="image-preview-item empty" onclick="document.getElementById('bannerImageInput').click()"><div style="text-align:center;"><span style="font-size:32px;">📁</span><p style="font-size:12px;color:#666;margin-top:5px;">上传Banner</p></div></div><input type="file" id="bannerImageInput" accept="image/*" onchange="handleBannerImageUpload(this)" style="display:none;">`;
    document.getElementById('bannerImageCount').textContent = '0';
    bannerUploadedFile = null;
    showToast('图片删除成功！');
  }
}

async function saveBannerConfig() {
  const formData = new FormData();
  formData.append('moduleName', 'banner');
  formData.append('enabled', document.getElementById('bannerEnabled').checked);
  formData.append('zhTitle', document.getElementById('bannerZhTitle').value);
  formData.append('enTitle', document.getElementById('bannerEnTitle').value);
  formData.append('zhSubtitle', document.getElementById('bannerZhSubtitle').value);
  formData.append('enSubtitle', document.getElementById('bannerEnSubtitle').value);
  formData.append('zhContent', document.getElementById('bannerZhContent').value);
  formData.append('enContent', document.getElementById('bannerEnContent').value);
  if (bannerUploadedFile) formData.append('image', bannerUploadedFile);
  const success = await saveModuleConfig('banner', formData);
  if (success) bannerUploadedFile = null;
}

function loadBannerConfig(banner) {
  document.getElementById('bannerEnabled').checked = banner.enabled !== false;
  document.getElementById('bannerZhTitle').value = banner.zhTitle || banner.title || '';
  document.getElementById('bannerEnTitle').value = banner.enTitle || '';
  document.getElementById('bannerZhSubtitle').value = banner.zhSubtitle || banner.subtitle || '';
  document.getElementById('bannerEnSubtitle').value = banner.enSubtitle || '';
  document.getElementById('bannerZhContent').value = banner.zhContent || banner.content || '';
  document.getElementById('bannerEnContent').value = banner.enContent || '';
  if (banner.imagePath) {
    document.getElementById('bannerImageCount').textContent = '1';
    document.getElementById('bannerImagesPreview').innerHTML = `<div class="image-preview-item"><img src="/uploads/${banner.imagePath}" alt="Banner"><div class="image-actions"><button class="action-btn replace" onclick="replaceBannerImage(event)">🔄</button><button class="action-btn delete" onclick="removeBannerImage()">🗑️</button></div></div><input type="file" id="bannerImageInput" accept="image/*" onchange="handleBannerImageUpload(this)" style="display:none;">`;
  }
}

function handleAboutImageUpload(input) {
  const file = input.files[0]; if (!file) return;
  const previewContainer = document.getElementById('aboutImagesPreview');
  const emptySlot = previewContainer.querySelector('.empty');
  if (!emptySlot) { showToast('关于我们模块仅允许上传1张图片', 'error'); return; }
  aboutUploadedFile = file;
  const reader = new FileReader();
  reader.onload = e => {
    emptySlot.innerHTML = `<img src="${e.target.result}" alt="About"><div class="image-actions"><button class="action-btn replace" onclick="replaceAboutImage(event)">🔄</button><button class="action-btn delete" onclick="removeAboutImage()">🗑️</button></div>`;
    emptySlot.classList.remove('empty');
    document.getElementById('aboutImageCount').textContent = '1';
    showToast('图片上传成功');
  };
  reader.readAsDataURL(file);
}

function replaceAboutImage(event) { event.stopPropagation(); document.getElementById('aboutImageInput').click(); }

async function removeAboutImage() {
  if (!confirm('确定要删除这张图片吗？')) return;
  const previewContainer = document.getElementById('aboutImagesPreview');
  const imageItem = previewContainer.querySelector('.image-preview-item:not(.empty)');
  if (imageItem) {
    const res = await apiFetch('/api/admin/modules/about/image', { method: 'DELETE' });
    if (!res.ok) { showToast('删除图片失败', 'error'); return; }
    previewContainer.innerHTML = `<div class="image-preview-item empty" onclick="document.getElementById('aboutImageInput').click()"><div style="text-align:center;"><span style="font-size:32px;">📁</span><p style="font-size:12px;color:#666;margin-top:5px;">上传图片</p></div></div><input type="file" id="aboutImageInput" accept="image/*" onchange="handleAboutImageUpload(this)" style="display:none;">`;
    document.getElementById('aboutImageCount').textContent = '0';
    aboutUploadedFile = null;
    showToast('图片删除成功！');
  }
}

async function saveAboutConfig() {
  const formData = new FormData();
  formData.append('moduleName', 'about');
  formData.append('enabled', document.getElementById('aboutEnabled').checked);
  formData.append('zhTitle', document.getElementById('aboutZhTitle').value);
  formData.append('enTitle', document.getElementById('aboutEnTitle').value);
  formData.append('zhSubtitle', document.getElementById('aboutZhSubtitle').value);
  formData.append('enSubtitle', document.getElementById('aboutEnSubtitle').value);
  formData.append('zhContent', document.getElementById('aboutZhContent').value);
  formData.append('enContent', document.getElementById('aboutEnContent').value);
  formData.append('extraData', JSON.stringify({ zhLeftTitle: document.getElementById('aboutZhLeftTitle').value, enLeftTitle: document.getElementById('aboutEnLeftTitle').value, leftTitle: document.getElementById('aboutZhLeftTitle').value }));
  if (aboutUploadedFile) formData.append('image', aboutUploadedFile);
  const success = await saveModuleConfig('about', formData);
  if (success) aboutUploadedFile = null;
}

function loadAboutConfig(about) {
  document.getElementById('aboutEnabled').checked = about.enabled !== false;
  document.getElementById('aboutZhTitle').value = about.zhTitle || about.title || '';
  document.getElementById('aboutEnTitle').value = about.enTitle || '';
  document.getElementById('aboutZhSubtitle').value = about.zhSubtitle || about.subtitle || '';
  document.getElementById('aboutEnSubtitle').value = about.enSubtitle || '';
  document.getElementById('aboutZhLeftTitle').value = about.zhLeftTitle || about.leftTitle || '';
  document.getElementById('aboutEnLeftTitle').value = about.enLeftTitle || '';
  document.getElementById('aboutZhContent').value = about.zhContent || about.content || '';
  document.getElementById('aboutEnContent').value = about.enContent || '';
  if (about.imagePath) {
    document.getElementById('aboutImageCount').textContent = '1';
    document.getElementById('aboutImagesPreview').innerHTML = `<div class="image-preview-item"><img src="/uploads/${about.imagePath}" alt="About"><div class="image-actions"><button class="action-btn replace" onclick="replaceAboutImage(event)">🔄</button><button class="action-btn delete" onclick="removeAboutImage()">🗑️</button></div></div><input type="file" id="aboutImageInput" accept="image/*" onchange="handleAboutImageUpload(this)" style="display:none;">`;
  }
}

async function loadProducts() {
  const res = await apiFetch('/api/admin/images?category=products');
  const images = (await res.json()).data || [];
  images.sort((a, b) => a.sortOrder - b.sortOrder);
  const list = document.getElementById('productsList');
  if (images.length === 0) {
    list.innerHTML = '<div style="text-align:center;color:#999;padding:50px;"><div style="font-size:48px;margin-bottom:15px;">📦</div><p>暂无产品数据，请添加产品</p></div>';
    return;
  }
  list.innerHTML = images.map(img => `<div class="product-item"><img src="/uploads/${img.filename}" alt="${img.description}"><h4>${img.description || '未命名产品'}</h4><p>${img.longDescription || '暂无描述'}</p><div class="actions"><button class="btn btn-primary" onclick="editProductItem(${img.id}, '${img.description || ''}', '${img.longDescription || ''}', '${img.filename}', ${img.sortOrder})">编辑</button><button class="btn btn-danger" onclick="deleteProductImage(${img.id})">删除</button></div></div>`).join('');
}

function addProductItem() {
  document.getElementById('productModalTitle').textContent = '添加产品';
  document.getElementById('productId').value = '';
  document.getElementById('productName').value = '';
  document.getElementById('productDescription').value = '';
  document.getElementById('productSortOrder').value = '0';
  document.getElementById('productImagePath').value = '';
  document.getElementById('productImageCount').textContent = '0';
  document.getElementById('productImagesPreview').innerHTML = `<div class="image-preview-item empty" onclick="document.getElementById('productImageInput').click()"><div style="text-align:center;"><span style="font-size:32px;">📁</span><p style="font-size:12px;color:#666;margin-top:5px;">添加图片</p></div></div><input type="file" id="productImageInput" accept="image/*" multiple onchange="handleProductImageUpload(this)" style="display:none;">`;
  productUploadedFiles = [];
  document.getElementById('productModal').classList.add('active');
}

function editProductItem(id, name, desc, filename, sortOrder) {
  document.getElementById('productModalTitle').textContent = '编辑产品';
  document.getElementById('productId').value = id;
  document.getElementById('productName').value = name;
  document.getElementById('productDescription').value = desc;
  document.getElementById('productSortOrder').value = sortOrder;
  document.getElementById('productImagePath').value = filename;
  productUploadedFiles = [];
  if (filename) {
    document.getElementById('productImageCount').textContent = '1';
    document.getElementById('productImagesPreview').innerHTML = `<div class="image-preview-item"><img src="/uploads/${filename}" alt="Product"><div class="image-actions"><button class="action-btn replace" onclick="replaceProductImage(event)">🔄</button><button class="action-btn delete" onclick="removeProductImage(this)">🗑️</button></div></div><div class="image-preview-item empty" onclick="document.getElementById('productImageInput').click()"><div style="text-align:center;"><span style="font-size:32px;">📁</span><p style="font-size:12px;color:#666;margin-top:5px;">添加图片</p></div></div><input type="file" id="productImageInput" accept="image/*" multiple onchange="handleProductImageUpload(this)" style="display:none;">`;
  } else {
    document.getElementById('productImageCount').textContent = '0';
    document.getElementById('productImagesPreview').innerHTML = `<div class="image-preview-item empty" onclick="document.getElementById('productImageInput').click()"><div style="text-align:center;"><span style="font-size:32px;">📁</span><p style="font-size:12px;color:#666;margin-top:5px;">添加图片</p></div></div><input type="file" id="productImageInput" accept="image/*" multiple onchange="handleProductImageUpload(this)" style="display:none;">`;
  }
  document.getElementById('productModal').classList.add('active');
}

function closeProductModal() { document.getElementById('productModal').classList.remove('active'); }

function handleProductImageUpload(input) {
  const files = input.files; if (!files || files.length === 0) return;
  const maxImages = 10;
  const previewContainer = document.getElementById('productImagesPreview');
  const currentCount = previewContainer.querySelectorAll('.image-preview-item:not(.empty)').length;
  const remainingSlots = maxImages - currentCount;
  const emptySlot = previewContainer.querySelector('.empty');
  if (files.length > remainingSlots) { showToast(`产品模块最多只能上传${maxImages}张图片，还可上传${remainingSlots}张`, 'error'); return; }
  for (let i = 0; i < files.length && i < remainingSlots; i++) {
    const file = files[i];
    productUploadedFiles.push(file);
    const reader = new FileReader();
    reader.onload = e => {
      const newItem = document.createElement('div');
      newItem.className = 'image-preview-item';
      newItem.innerHTML = `<img src="${e.target.result}" alt="Product"><div class="image-actions"><button class="action-btn replace" onclick="replaceProductImage(event)">🔄</button><button class="action-btn delete" onclick="removeProductImage(this)">🗑️</button></div>`;
      previewContainer.insertBefore(newItem, emptySlot);
      document.getElementById('productImageCount').textContent = currentCount + i + 1;
    };
    reader.readAsDataURL(file);
  }
  showToast(`成功上传${Math.min(files.length, remainingSlots)}张图片`);
}

function replaceProductImage(event) { event.stopPropagation(); document.getElementById('productImageInput').click(); }

function removeProductImage(btn) {
  const previewContainer = document.getElementById('productImagesPreview');
  const item = btn.closest('.image-preview-item');
  if (item) {
    const index = Array.from(previewContainer.children).indexOf(item);
    if (index > -1 && index < productUploadedFiles.length) productUploadedFiles.splice(index, 1);
    item.remove();
    document.getElementById('productImageCount').textContent = previewContainer.querySelectorAll('.image-preview-item:not(.empty)').length;
  }
}

async function saveProductItem() {
  const id = document.getElementById('productId').value;
  const formData = new FormData();
  formData.append('category', 'products');
  formData.append('description', document.getElementById('productName').value);
  formData.append('longDescription', document.getElementById('productDescription').value);
  formData.append('sortOrder', document.getElementById('productSortOrder').value);
  if (productUploadedFiles.length > 0) formData.append('image', productUploadedFiles[0]);
  else formData.append('imagePath', document.getElementById('productImagePath').value);
  const isValidProductId = id && id !== 'undefined' && id.trim() !== '';
  const url = isValidProductId ? `/api/admin/images/${id}` : '/api/admin/images';
  const res = await apiFetch(url, { method: isValidProductId ? 'PUT' : 'POST', body: formData });
  if (res.ok) { closeProductModal(); loadProducts(); showToast('产品保存成功！'); }
  else showToast('保存失败', 'error');
}

async function deleteProductImage(id) {
  if (!confirm('确定删除此产品？')) return;
  const res = await apiFetch(`/api/admin/images/${id}`, { method: 'DELETE' });
  if (res.ok) { loadProducts(); showToast('产品删除成功！'); }
}

function loadProductsConfig(products) {
  document.getElementById('productsEnabled').checked = products.enabled !== false;
  document.getElementById('productsTitle').value = products.title || '';
  document.getElementById('productsCount').value = products.showCount || 6;
}

async function saveProductsModuleConfig() {
  const formData = new FormData();
  formData.append('moduleName', 'products');
  formData.append('enabled', document.getElementById('productsEnabled').checked);
  formData.append('title', document.getElementById('productsTitle').value);
  formData.append('extraData', JSON.stringify({ showCount: parseInt(document.getElementById('productsCount').value) }));
  const res = await apiFetch('/api/admin/modules', { method: 'POST', body: formData });
  if (res.ok) showToast('配置保存成功！');
  else { const err = await res.json(); showToast('保存失败: ' + (err.error || '未知错误'), 'error'); }
}

function loadFactoryConfig(factory) {
  document.getElementById('factoryEnabled').checked = factory.enabled !== false;
  document.getElementById('factoryTitle').value = factory.title || '';
  document.getElementById('factorySubtitle').value = factory.subtitle || '';
}

async function saveFactoryModuleConfig() {
  const formData = new FormData();
  formData.append('moduleName', 'factory');
  formData.append('enabled', document.getElementById('factoryEnabled').checked);
  formData.append('title', document.getElementById('factoryTitle').value);
  formData.append('subtitle', document.getElementById('factorySubtitle').value);
  const res = await apiFetch('/api/admin/modules', { method: 'POST', body: formData });
  if (res.ok) showToast('配置保存成功！');
  else { const err = await res.json(); showToast('保存失败: ' + (err.error || '未知错误'), 'error'); }
}

function loadAdvantageConfig(advantage) {
  document.getElementById('advantageEnabled').checked = advantage.enabled !== false;
  document.getElementById('advantageTitle').value = advantage.title || '';
  document.getElementById('advantageSubtitle').value = advantage.subtitle || '';
}

async function saveAdvantageModuleConfig() {
  const formData = new FormData();
  formData.append('moduleName', 'advantage');
  formData.append('enabled', document.getElementById('advantageEnabled').checked);
  formData.append('title', document.getElementById('advantageTitle').value);
  formData.append('subtitle', document.getElementById('advantageSubtitle').value);
  const res = await apiFetch('/api/admin/modules', { method: 'POST', body: formData });
  if (res.ok) showToast('配置保存成功！');
  else { const err = await res.json(); showToast('保存失败: ' + (err.error || '未知错误'), 'error'); }
}

function addFactoryItem() {
  document.getElementById('contentModalTitle').textContent = '添加工厂优势项目';
  document.getElementById('contentSectionType').value = 'factory';
  document.getElementById('contentItemId').value = '';
  document.getElementById('contentZhTitle').value = '';
  document.getElementById('contentEnTitle').value = '';
  document.getElementById('contentZhDescription').value = '';
  document.getElementById('contentEnDescription').value = '';
  document.getElementById('contentIcon').value = '';
  document.getElementById('contentSortOrder').value = '0';
  document.getElementById('contentImagePath').value = '';
  document.getElementById('contentImageCount').textContent = '0';
  document.getElementById('contentImagesPreview').innerHTML = `<div class="image-preview-item empty" onclick="document.getElementById('contentImageInput').click()"><div style="text-align:center;"><span style="font-size:32px;">📁</span><p style="font-size:12px;color:#666;margin-top:5px;">上传图标（可选）</p></div></div><input type="file" id="contentImageInput" accept="image/*" onchange="handleContentImageUpload(this)" style="display:none;">`;
  contentUploadedFile = null;
  document.getElementById('contentModal').classList.add('active');
}

function addAdvantageItem() {
  document.getElementById('contentModalTitle').textContent = '添加核心优势项目';
  document.getElementById('contentSectionType').value = 'advantage';
  document.getElementById('contentItemId').value = '';
  document.getElementById('contentZhTitle').value = '';
  document.getElementById('contentEnTitle').value = '';
  document.getElementById('contentZhDescription').value = '';
  document.getElementById('contentEnDescription').value = '';
  document.getElementById('contentIcon').value = '';
  document.getElementById('contentSortOrder').value = '0';
  document.getElementById('contentImagePath').value = '';
  document.getElementById('contentImageCount').textContent = '0';
  document.getElementById('contentImagesPreview').innerHTML = `<div class="image-preview-item empty" onclick="document.getElementById('contentImageInput').click()"><div style="text-align:center;"><span style="font-size:32px;">📁</span><p style="font-size:12px;color:#666;margin-top:5px;">上传图标（可选）</p></div></div><input type="file" id="contentImageInput" accept="image/*" onchange="handleContentImageUpload(this)" style="display:none;">`;
  contentUploadedFile = null;
  document.getElementById('contentModal').classList.add('active');
}

function closeContentModal() { document.getElementById('contentModal').classList.remove('active'); }

function handleContentImageUpload(input) {
  const file = input.files[0]; if (!file) return;
  const previewContainer = document.getElementById('contentImagesPreview');
  const emptySlot = previewContainer.querySelector('.empty');
  if (!emptySlot) { showToast('此模块仅允许上传1张图片', 'error'); return; }
  contentUploadedFile = file;
  const reader = new FileReader();
  reader.onload = e => {
    emptySlot.innerHTML = `<img src="${e.target.result}" alt="Icon"><div class="image-actions"><button class="action-btn replace" onclick="replaceContentImage(event)">🔄</button><button class="action-btn delete" onclick="removeContentImage()">🗑️</button></div>`;
    emptySlot.classList.remove('empty');
    document.getElementById('contentImageCount').textContent = '1';
    showToast('图片上传成功');
  };
  reader.readAsDataURL(file);
}

function replaceContentImage(event) { event.stopPropagation(); document.getElementById('contentImageInput').click(); }

function removeContentImage() {
  const previewContainer = document.getElementById('contentImagesPreview');
  const imageItem = previewContainer.querySelector('.image-preview-item:not(.empty)');
  if (imageItem) {
    previewContainer.innerHTML = `<div class="image-preview-item empty" onclick="document.getElementById('contentImageInput').click()"><div style="text-align:center;"><span style="font-size:32px;">📁</span><p style="font-size:12px;color:#666;margin-top:5px;">上传图标（可选）</p></div></div><input type="file" id="contentImageInput" accept="image/*" onchange="handleContentImageUpload(this)" style="display:none;">`;
    document.getElementById('contentImageCount').textContent = '0';
    contentUploadedFile = null;
  }
}

async function saveContentItem() {
  const id = document.getElementById('contentItemId').value;
  const sectionType = document.getElementById('contentSectionType').value;
  const formData = new FormData();
  formData.append('section', sectionType);
  formData.append('zhTitle', document.getElementById('contentZhTitle').value);
  formData.append('enTitle', document.getElementById('contentEnTitle').value);
  formData.append('zhDescription', document.getElementById('contentZhDescription').value);
  formData.append('enDescription', document.getElementById('contentEnDescription').value);
  formData.append('icon', document.getElementById('contentIcon').value);
  formData.append('sortOrder', document.getElementById('contentSortOrder').value);
  if (contentUploadedFile) formData.append('image', contentUploadedFile);
  else formData.append('imagePath', document.getElementById('contentImagePath').value);
  const isValidId = id && id !== 'undefined' && id.trim() !== '';
  const url = isValidId ? `/api/admin/content/${id}` : '/api/admin/content';
  const res = await apiFetch(url, { method: isValidId ? 'PUT' : 'POST', body: formData });
  if (res.ok) { closeContentModal(); loadFactoryItems(); loadAdvantageItems(); showToast('保存成功！'); }
  else showToast('保存失败', 'error');
}

async function loadFactoryItems() {
  const res = await apiFetch('/api/admin/content?section=factory');
  const items = (await res.json()).data || [];
  items.sort((a, b) => a.sortOrder - b.sortOrder);
  const list = document.getElementById('factoryList');
  if (items.length === 0) {
    list.innerHTML = '<div style="text-align:center;color:#999;padding:50px;"><div style="font-size:48px;margin-bottom:15px;">🏭</div><p>暂无工厂优势项目，请添加</p></div>';
    return;
  }
  list.innerHTML = items.map(item => `<div class="content-item"><div class="icon">${item.icon || '🏭'}</div><div class="info"><h4>${item.zhTitle || item.title || '未命名'}</h4><p>${item.zhDescription || '暂无描述'}</p></div><div class="actions"><button class="btn btn-primary" onclick="editContentItem(${item.id}, 'factory', '${item.zhTitle || ''}', '${item.enTitle || ''}', '${item.zhDescription || ''}', '${item.enDescription || ''}', '${item.icon || ''}', ${item.sortOrder})">编辑</button><button class="btn btn-danger" onclick="deleteContentItem(${item.id}, 'factory')">删除</button></div></div>`).join('');
}

async function loadAdvantageItems() {
  const res = await apiFetch('/api/admin/content?section=advantage');
  const items = (await res.json()).data || [];
  items.sort((a, b) => a.sortOrder - b.sortOrder);
  const list = document.getElementById('advantageList');
  if (items.length === 0) {
    list.innerHTML = '<div style="text-align:center;color:#999;padding:50px;"><div style="font-size:48px;margin-bottom:15px;">⭐</div><p>暂无核心优势项目，请添加</p></div>';
    return;
  }
  list.innerHTML = items.map(item => `<div class="content-item"><div class="icon">${item.icon || '⭐'}</div><div class="info"><h4>${item.zhTitle || item.title || '未命名'}</h4><p>${item.zhDescription || '暂无描述'}</p></div><div class="actions"><button class="btn btn-primary" onclick="editContentItem(${item.id}, 'advantage', '${item.zhTitle || ''}', '${item.enTitle || ''}', '${item.zhDescription || ''}', '${item.enDescription || ''}', '${item.icon || ''}', ${item.sortOrder})">编辑</button><button class="btn btn-danger" onclick="deleteContentItem(${item.id}, 'advantage')">删除</button></div></div>`).join('');
}

function editContentItem(id, type, zhTitle, enTitle, zhDesc, enDesc, icon, sortOrder) {
  document.getElementById('contentModalTitle').textContent = type === 'factory' ? '编辑工厂优势项目' : '编辑核心优势项目';
  document.getElementById('contentSectionType').value = type;
  document.getElementById('contentItemId').value = id;
  document.getElementById('contentZhTitle').value = zhTitle;
  document.getElementById('contentEnTitle').value = enTitle;
  document.getElementById('contentZhDescription').value = zhDesc;
  document.getElementById('contentEnDescription').value = enDesc;
  document.getElementById('contentIcon').value = icon;
  document.getElementById('contentSortOrder').value = sortOrder;
  document.getElementById('contentModal').classList.add('active');
}

async function deleteContentItem(id, type) {
  if (!confirm('确定删除此项？')) return;
  const res = await apiFetch(`/api/admin/content/${id}`, { method: 'DELETE' });
  if (res.ok) { loadFactoryItems(); loadAdvantageItems(); showToast('删除成功！'); }
}

function loadEventsConfig(events) {
  document.getElementById('eventsEnabled').checked = events.enabled !== false;
  document.getElementById('eventsZhName').value = events.zhName || events.name || '';
  document.getElementById('eventsEnName').value = events.enName || '';
  document.getElementById('eventsBooth').value = events.booth || '';
  document.getElementById('eventsStartDate').value = events.startDate || '';
  document.getElementById('eventsEndDate').value = events.endDate || '';
  document.getElementById('eventsZhLocation').value = events.zhLocation || events.location || '';
  document.getElementById('eventsEnLocation').value = events.enLocation || '';
  document.getElementById('eventsZhDescription').value = events.zhDescription || events.description || '';
  document.getElementById('eventsEnDescription').value = events.enDescription || '';
  document.getElementById('eventsIcon').value = events.icon || '';
  document.getElementById('eventsZhLeftTitle').value = events.zhLeftTitle || events.leftTitle || '';
  document.getElementById('eventsEnLeftTitle').value = events.enLeftTitle || '';
  document.getElementById('eventsZhLeftSubtitle').value = events.zhLeftSubtitle || events.leftSubtitle || '';
  document.getElementById('eventsEnLeftSubtitle').value = events.enLeftSubtitle || '';
}

function validateEventDates() {
  const startDateInput = document.getElementById('eventsStartDate');
  const endDateInput = document.getElementById('eventsEndDate');
  const startDate = startDateInput.value;
  const endDate = endDateInput.value;
  
  clearDateValidationErrors();
  
  if (!startDate && endDate) {
    showDateError(endDateInput, '请先选择展会的开始日期');
    return false;
  }
  
  if (startDate && endDate) {
    const start = new Date(startDate);
    const end = new Date(endDate);
    
    if (end < start) {
      showDateError(endDateInput, '结束日期不能早于开始日期');
      return false;
    }
  }
  
  return true;
}

function showDateError(input, message) {
  input.style.borderColor = '#dc3545';
  input.style.backgroundColor = '#fff8f8';
  
  let errorEl = input.parentNode.querySelector('.date-error');
  if (!errorEl) {
    errorEl = document.createElement('div');
    errorEl.className = 'date-error';
    errorEl.style.cssText = 'color:#dc3545;font-size:12px;margin-top:5px;display:flex;align-items:center;gap:5px;';
    input.parentNode.appendChild(errorEl);
  }
  errorEl.innerHTML = `⚠️ ${message}`;
}

function clearDateValidationErrors() {
  ['eventsStartDate', 'eventsEndDate'].forEach(id => {
    const input = document.getElementById(id);
    if (input) {
      input.style.borderColor = '';
      input.style.backgroundColor = '';
      const errorEl = input.parentNode.querySelector('.date-error');
      if (errorEl) errorEl.remove();
    }
  });
}

function initEventsDateValidation() {
  const startDateInput = document.getElementById('eventsStartDate');
  const endDateInput = document.getElementById('eventsEndDate');
  
  if (startDateInput) {
    startDateInput.addEventListener('change', function() {
      if (this.value && endDateInput.value) {
        validateEventDates();
      } else {
        clearDateValidationErrors();
      }
    });
  }
  
  if (endDateInput) {
    endDateInput.addEventListener('change', function() {
      validateEventDates();
    });
  }
}

async function saveEventsConfig() {
  if (!validateEventDates()) {
    showToast('请修正日期错误后再保存', 'error');
    return false;
  }
  
  const formData = new FormData();
  formData.append('moduleName', 'events');
  formData.append('enabled', document.getElementById('eventsEnabled').checked);
  formData.append('zhName', document.getElementById('eventsZhName').value);
  formData.append('enName', document.getElementById('eventsEnName').value);
  formData.append('booth', document.getElementById('eventsBooth').value);
  formData.append('startDate', document.getElementById('eventsStartDate').value);
  formData.append('endDate', document.getElementById('eventsEndDate').value);
  formData.append('zhLocation', document.getElementById('eventsZhLocation').value);
  formData.append('enLocation', document.getElementById('eventsEnLocation').value);
  formData.append('zhDescription', document.getElementById('eventsZhDescription').value);
  formData.append('enDescription', document.getElementById('eventsEnDescription').value);
  formData.append('extraData', JSON.stringify({
    icon: document.getElementById('eventsIcon').value,
    zhLeftTitle: document.getElementById('eventsZhLeftTitle').value,
    enLeftTitle: document.getElementById('eventsEnLeftTitle').value,
    zhLeftSubtitle: document.getElementById('eventsZhLeftSubtitle').value,
    enLeftSubtitle: document.getElementById('eventsEnLeftSubtitle').value
  }));
  const success = await saveModuleConfig('events', formData);
  if (success) {
    clearDateValidationErrors();
  }
}

function loadContactConfig(contact) {
  document.getElementById('contactEnabled').checked = contact.enabled !== false;
  document.getElementById('contactEmail').value = contact.email || '';
  document.getElementById('contactPhone').value = contact.phone || '';
  document.getElementById('contactWhatsApp').value = contact.whatsapp || contact.whatsApp || '';
  document.getElementById('contactAddress').value = contact.address || '';
}

async function saveContactConfig() {
  const formData = new FormData();
  formData.append('moduleName', 'contact');
  formData.append('enabled', document.getElementById('contactEnabled').checked);
  formData.append('email', document.getElementById('contactEmail').value);
  formData.append('phone', document.getElementById('contactPhone').value);
  formData.append('whatsapp', document.getElementById('contactWhatsApp').value);
  formData.append('address', document.getElementById('contactAddress').value);
  const success = await saveModuleConfig('contact', formData);
}

function loadSystemConfig(brand, stats) {
  if (brand) {
    document.getElementById('brandName').value = brand.name || '';
    document.getElementById('brandSuffixColor').value = brand.suffixColor || '#0a5cad';
  }
  if (stats) {
    document.getElementById('statsYearsInput').value = stats.years || '15+';
    document.getElementById('statsProductsInput').value = stats.products || '200+';
    document.getElementById('statsCountriesInput').value = stats.countries || '80+';
    document.getElementById('statsClientsInput').value = stats.clients || '1000+';
    document.getElementById('statsYears').textContent = stats.years || '15+';
    document.getElementById('statsProducts').textContent = stats.products || '200+';
    document.getElementById('statsCountries').textContent = stats.countries || '80+';
    document.getElementById('statsClients').textContent = stats.clients || '1000+';
  }
}

async function saveBrandConfig() {
  const config = { name: document.getElementById('brandName').value, suffixColor: document.getElementById('brandSuffixColor').value };
  await savePageConfig('brand', config);
}

async function saveStatsConfig() {
  const config = {
    years: document.getElementById('statsYearsInput').value,
    products: document.getElementById('statsProductsInput').value,
    countries: document.getElementById('statsCountriesInput').value,
    clients: document.getElementById('statsClientsInput').value
  };
  const success = await savePageConfig('stats', config);
  if (success) {
    document.getElementById('statsYears').textContent = config.years;
    document.getElementById('statsProducts').textContent = config.products;
    document.getElementById('statsCountries').textContent = config.countries;
    document.getElementById('statsClients').textContent = config.clients;
  }
}

function addLangText() {
  document.getElementById('langModalTitle').textContent = '添加文案';
  document.getElementById('langTextId').value = '';
  document.getElementById('langKey').value = '';
  document.getElementById('langModule').value = 'nav';
  document.getElementById('langEnText').value = '';
  document.getElementById('langZhText').value = '';
  document.getElementById('langDescription').value = '';
  document.getElementById('langModal').classList.add('active');
}

function closeLangModal() { document.getElementById('langModal').classList.remove('active'); }

async function saveLangText() {
  const id = document.getElementById('langTextId').value;
  const data = {
    key: document.getElementById('langKey').value,
    module: document.getElementById('langModule').value,
    enText: document.getElementById('langEnText').value,
    zhText: document.getElementById('langZhText').value,
    description: document.getElementById('langDescription').value
  };
  const url = id && id !== 'undefined' && id.trim() !== '' ? `/api/admin/lang/${id}` : '/api/admin/lang';
  const res = await apiFetch(url, { method: id && id !== 'undefined' && id.trim() !== '' ? 'PUT' : 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(data) });
  if (res.ok) { closeLangModal(); loadLangTexts(); showToast('文案保存成功！'); }
  else showToast('保存失败', 'error');
}

async function loadLangTexts() {
  const module = document.getElementById('langModuleSelect').value;
  const url = module ? `/api/admin/lang?module=${module}` : '/api/admin/lang';
  const res = await apiFetch(url);
  const texts = (await res.json()).data || [];
  const list = document.getElementById('langList');
  if (texts.length === 0) {
    list.innerHTML = '<div style="text-align:center;color:#999;padding:50px;"><div style="font-size:48px;margin-bottom:15px;">📝</div><p>暂无文案数据，请添加文案</p></div>';
    return;
  }
  list.innerHTML = texts.map(t => `<div class="content-item"><div class="icon">📝</div><div class="info"><h4>${t.key}</h4><p>${t.zhText || t.enText || '暂无内容'}</p></div><div class="actions"><button class="btn btn-primary" onclick="editLangText(${t.id}, '${t.key}', '${t.module}', '${t.enText || ''}', '${t.zhText || ''}', '${t.description || ''}')">编辑</button><button class="btn btn-danger" onclick="deleteLangText(${t.id})">删除</button></div></div>`).join('');
}

function editLangText(id, key, module, enText, zhText, description) {
  document.getElementById('langModalTitle').textContent = '编辑文案';
  document.getElementById('langTextId').value = id;
  document.getElementById('langKey').value = key;
  document.getElementById('langModule').value = module;
  document.getElementById('langEnText').value = enText;
  document.getElementById('langZhText').value = zhText;
  document.getElementById('langDescription').value = description;
  document.getElementById('langModal').classList.add('active');
}

async function deleteLangText(id) {
  if (!confirm('确定删除此文案？')) return;
  const res = await apiFetch(`/api/admin/lang/${id}`, { method: 'DELETE' });
  if (res.ok) { loadLangTexts(); showToast('删除成功！'); }
}

async function loadSiteSettings() {
  try {
    const res = await apiFetch('/api/admin/site-settings');
    if (res.ok) {
      const data = (await res.json()).data || {};
      document.getElementById('zhSiteTitle').value = data.zhSiteTitle || '';
      document.getElementById('enSiteTitle').value = data.enSiteTitle || '';
      document.getElementById('zhSiteLogo').value = data.zhSiteLogo || '';
      document.getElementById('enSiteLogo').value = data.enSiteLogo || '';
      document.getElementById('siteLogoColor').value = data.siteLogoColor || '#06a499';
    }
  } catch(e) {
    console.error('Error loading site settings:', e);
  }
}

async function saveSiteSettings() {
  const zhTitle = document.getElementById('zhSiteTitle').value;
  const enTitle = document.getElementById('enSiteTitle').value;
  const zhLogo = document.getElementById('zhSiteLogo').value;
  const enLogo = document.getElementById('enSiteLogo').value;
  const logoColor = document.getElementById('siteLogoColor').value;
  
  if (!zhTitle && !enTitle && !zhLogo && !enLogo) {
    showToast('至少填写一个字段', 'error');
    return;
  }
  
  const formData = new FormData();
  formData.append('zhSiteTitle', zhTitle);
  formData.append('enSiteTitle', enTitle);
  formData.append('zhSiteLogo', zhLogo);
  formData.append('enSiteLogo', enLogo);
  formData.append('siteLogoColor', logoColor);
  
  const res = await apiFetch('/api/admin/site-settings', { method: 'POST', body: formData });
  if (res.ok) {
    showToast('网站配置保存成功！');
    // 重新加载设置以回显数据
    await loadSiteSettings();
  } else {
    const err = await res.json();
    showToast('保存失败：' + (err.error || '未知错误'), 'error');
  }
}

function showSection(sectionId) {
  document.querySelectorAll('.nav-item').forEach(item => item.classList.remove('active'));
  document.querySelector(`[data-section="${sectionId}"]`).classList.add('active');
  document.querySelectorAll('.section').forEach(s => s.classList.remove('active'));
  document.getElementById(`${sectionId}Section`).classList.add('active');
  const titles = { banner: '首页 Banner 配置', about: '关于我们配置', products: '产品展示配置', factory: '工厂实力配置', advantage: '核心优势配置', events: '展会信息配置', contact: '联系我们配置', system: '系统设置', lang: '多语言配置' };
  document.getElementById('pageTitle').textContent = titles[sectionId] || '配置';
  if (sectionId === 'products') loadProducts();
  if (sectionId === 'factory') loadFactoryItems();
  if (sectionId === 'advantage') loadAdvantageItems();
  if (sectionId === 'lang') loadLangTexts();
  if (sectionId === 'events') initEventsDateValidation();
  if (sectionId === 'contact') loadContactSubmissions();
  window.location.hash = sectionId;
}

function initApp() {
  if (token) {
    document.getElementById('loginPage').style.display = 'none';
    document.getElementById('adminPage').style.display = 'flex';
    loadAllConfigs();
    const hash = window.location.hash.replace('#', '');
    const validSections = ['banner', 'about', 'products', 'factory', 'advantage', 'events', 'contact', 'system', 'lang'];
    if (hash && validSections.includes(hash)) {
      showSection(hash);
    }
  }
}

// 联系表单提交数据管理
let contactSubmissionsCurrentPage = 1;
let contactSubmissionsTotalPage = 1;

// 加载联系表单提交数据
async function loadContactSubmissions(page = 1) {
  try {
    const res = await apiFetch(`/api/admin/contact-submissions?page=${page}&page_size=10`);
    const responseData = await res.json();
    
    if (responseData.code === 200 && responseData.data) {
      const { list, total, page: currentPage, total_page } = responseData.data;
      
      contactSubmissionsCurrentPage = currentPage;
      contactSubmissionsTotalPage = total_page;
      
      // 更新分页信息
      document.getElementById('contactSubmissionsTotal').textContent = total;
      document.getElementById('contactSubmissionsCurrentPage').textContent = currentPage;
      document.getElementById('contactSubmissionsTotalPage').textContent = total_page;
      
      // 显示分页控件
      const paginationEl = document.getElementById('contactSubmissionsPagination');
      if (total > 0) {
        paginationEl.style.display = 'block';
      } else {
        paginationEl.style.display = 'none';
      }
      
      // 更新按钮状态
      document.getElementById('contactPrevBtn').disabled = currentPage <= 1;
      document.getElementById('contactNextBtn').disabled = currentPage >= total_page;
      
      // 渲染列表
      const listEl = document.getElementById('contactSubmissionsList');
      if (!list || list.length === 0) {
        listEl.innerHTML = `
          <div style="text-align:center;color:#999;padding:50px;">
            <div style="font-size:48px;margin-bottom:15px;">📬</div>
            <p>暂无提交记录</p>
          </div>
        `;
        return;
      }
      
      let html = `
        <table style="width:100%;border-collapse:collapse;">
          <thead>
            <tr style="background:#f5f5f5;">
              <th style="padding:12px;text-align:left;border-bottom:2px solid #ddd;">姓名</th>
              <th style="padding:12px;text-align:left;border-bottom:2px solid #ddd;">邮箱</th>
              <th style="padding:12px;text-align:left;border-bottom:2px solid #ddd;">公司/国家</th>
              <th style="padding:12px;text-align:left;border-bottom:2px solid #ddd;">询盘内容</th>
              <th style="padding:12px;text-align:left;border-bottom:2px solid #ddd;">提交时间</th>
              <th style="padding:12px;text-align:center;border-bottom:2px solid #ddd;">操作</th>
            </tr>
          </thead>
          <tbody>
      `;
      
      list.forEach(item => {
        const isRead = item.isRead ? '是' : '否';
        const readStyle = item.isRead ? 'color:#999;' : 'color:#0a5cad;font-weight:bold;';
        const date = new Date(item.createdAt).toLocaleString('zh-CN');
        
        html += `
          <tr style="border-bottom:1px solid #f0f0f0;">
            <td style="padding:12px;${readStyle}">${escapeHtml(item.name)}</td>
            <td style="padding:12px;">${escapeHtml(item.email)}</td>
            <td style="padding:12px;">${escapeHtml(item.company || '-')}</td>
            <td style="padding:12px;max-width:300px;">
              <div style="overflow:hidden;text-overflow:ellipsis;white-space:nowrap;" title="${escapeHtml(item.inquiry || '')}">
                ${escapeHtml(item.inquiry || '-')}
              </div>
            </td>
            <td style="padding:12px;color:#666;">${date}</td>
            <td style="padding:12px;text-align:center;">
              <button class="btn btn-secondary" style="padding:5px 10px;font-size:12px;" onclick="deleteContactSubmission(${item.id})">删除</button>
            </td>
          </tr>
        `;
      });
      
      html += '</tbody></table>';
      listEl.innerHTML = html;
    }
  } catch (e) {
    console.error('Error loading contact submissions:', e);
    showToast('加载失败', 'error');
  }
}

// 切换分页
function changeContactSubmissionsPage(delta) {
  const newPage = contactSubmissionsCurrentPage + delta;
  if (newPage >= 1 && newPage <= contactSubmissionsTotalPage) {
    loadContactSubmissions(newPage);
  }
}

// 删除联系表单提交记录
async function deleteContactSubmission(id) {
  if (!confirm('确定要删除这条提交记录吗？')) return;
  
  try {
    const res = await apiFetch(`/api/admin/contact-submissions/${id}`, { method: 'DELETE' });
    const responseData = await res.json();
    
    if (responseData.code === 200) {
      showToast('删除成功');
      loadContactSubmissions(contactSubmissionsCurrentPage);
    } else {
      showToast('删除失败：' + (responseData.message || '未知错误'), 'error');
    }
  } catch (e) {
    console.error('Error deleting contact submission:', e);
    showToast('删除失败', 'error');
  }
}

// HTML 转义函数
function escapeHtml(text) {
  if (!text) return '';
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}

document.addEventListener('DOMContentLoaded', () => {
  initApp();
  window.addEventListener('hashchange', () => {
    const hash = window.location.hash.replace('#', '');
    const validSections = ['banner', 'about', 'products', 'factory', 'advantage', 'events', 'contact', 'system', 'lang'];
    if (hash && validSections.includes(hash)) {
      showSection(hash);
    }
  });
});
