const express = require('express');
const sqlite3 = require('sqlite3').verbose();
const multer = require('multer');
const bcrypt = require('bcryptjs');
const session = require('express-session');
const cors = require('cors');
const path = require('path');
const fs = require('fs');

const app = express();
const PORT = 3000;

// 配置
app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use('/uploads', express.static('uploads'));
app.use(express.static('.'));

// 会话管理
app.use(session({
    secret: 'medical-cms-secret-key',
    resave: false,
    saveUninitialized: false,
    cookie: { maxAge: 86400000 }
}));

// 确保上传目录存在
if (!fs.existsSync('uploads')) {
    fs.mkdirSync('uploads');
}

// 图片上传配置
const storage = multer.diskStorage({
    destination: (req, file, cb) => {
        cb(null, 'uploads/');
    },
    filename: (req, file, cb) => {
        const uniqueSuffix = Date.now() + '-' + Math.round(Math.random() * 1E9);
        cb(null, uniqueSuffix + path.extname(file.originalname));
    }
});

const upload = multer({ 
    storage: storage,
    fileFilter: (req, file, cb) => {
        const allowedTypes = /jpeg|jpg|png|gif|webp/;
        const extname = allowedTypes.test(path.extname(file.originalname).toLowerCase());
        const mimetype = allowedTypes.test(file.mimetype);

        if (mimetype && extname) {
            return cb(null, true);
        } else {
            cb(new Error('只允许上传图片文件'));
        }
    }
});

// 数据库初始化
const db = new sqlite3.Database('./cms.db', (err) => {
    if (err) {
        console.error(err.message);
    }
    console.log('已连接到SQLite数据库');
    initDatabase();
});

// 初始化数据库表
function initDatabase() {
    // 用户表
    db.run(`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL,
        role TEXT NOT NULL DEFAULT 'admin',
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`);

    // 图片表
    db.run(`CREATE TABLE IF NOT EXISTS images (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        filename TEXT NOT NULL,
        original_name TEXT NOT NULL,
        path TEXT NOT NULL,
        url TEXT NOT NULL,
        category TEXT NOT NULL,
        description TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`);

    // 页面配置表
    db.run(`CREATE TABLE IF NOT EXISTS page_configs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        page_name TEXT NOT NULL,
        section_name TEXT NOT NULL,
        config_key TEXT NOT NULL,
        config_value TEXT,
        config_type TEXT DEFAULT 'text',
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`);

    // 插入默认管理员用户
    const defaultPassword = bcrypt.hashSync('admin123', 10);
    db.run(`INSERT OR IGNORE INTO users (username, password, role) VALUES (?, ?, ?)`, 
        ['admin', defaultPassword, 'admin'], 
        (err) => {
            if (err) {
                console.error(err.message);
            } else {
                console.log('默认管理员用户已创建 (用户名: admin, 密码: admin123)');
            }
        });

    // 初始化默认配置
    initDefaultConfigs();
}

// 初始化默认配置
function initDefaultConfigs() {
    const defaultConfigs = [
        { page: 'home', section: 'banner', key: 'title', value: 'Professional Medical Device Manufacturer & Supplier' },
        { page: 'home', section: 'banner', key: 'subtitle', value: 'CE, FDA, ISO Certified Medical Devices. We provide high-quality hospital equipment, surgical supplies and disposable medical products with OEM/ODM global wholesale service.' },
        { page: 'home', section: 'about', key: 'title', value: 'About Our Company' },
        { page: 'home', section: 'about', key: 'subtitle', value: 'Reliable & Professional Medical Device Manufacturer from China' },
        { page: 'home', section: 'products', key: 'title', value: 'Our Main Products' },
        { page: 'home', section: 'products', key: 'subtitle', value: 'High-quality, sterile, safe medical devices for global hospital & medical institution use' },
        { page: 'home', section: 'factory', key: 'title', value: 'Our Factory Strength' },
        { page: 'home', section: 'factory', key: 'subtitle', value: 'Standard Dust-free Workshop, Strict Medical Grade Production & QC System' },
        { page: 'home', section: 'advantage', key: 'title', value: 'Why Choose Us' },
        { page: 'home', section: 'advantage', key: 'subtitle', value: 'Professional, Safe, Reliable Global Medical Device Supplier' },
        { page: 'home', section: 'contact', key: 'title', value: 'Contact Us' },
        { page: 'home', section: 'contact', key: 'subtitle', value: 'Get the latest medical device wholesale price and customized solution' },
    ];

    defaultConfigs.forEach(config => {
        db.run(`INSERT OR IGNORE INTO page_configs (page_name, section_name, config_key, config_value) VALUES (?, ?, ?, ?)`,
            [config.page, config.section, config.key, config.value]);
    });
}

// 权限中间件
const requireAuth = (req, res, next) => {
    if (req.session.user) {
        next();
    } else {
        res.status(401).json({ error: '需要登录' });
    }
};

// ==================== 用户认证 API ====================

// 登录
app.post('/api/login', (req, res) => {
    const { username, password } = req.body;

    db.get(`SELECT * FROM users WHERE username = ?`, [username], (err, user) => {
        if (err) {
            return res.status(500).json({ error: err.message });
        }
        if (!user) {
            return res.status(401).json({ error: '用户名或密码错误' });
        }

        const passwordMatch = bcrypt.compareSync(password, user.password);
        if (!passwordMatch) {
            return res.status(401).json({ error: '用户名或密码错误' });
        }

        req.session.user = { id: user.id, username: user.username, role: user.role };
        res.json({ success: true, user: { id: user.id, username: user.username, role: user.role } });
    });
});

// 登出
app.post('/api/logout', (req, res) => {
    req.session.destroy();
    res.json({ success: true });
});

// 获取当前用户
app.get('/api/user', (req, res) => {
    if (req.session.user) {
        res.json({ user: req.session.user });
    } else {
        res.status(401).json({ error: '未登录' });
    }
});

// ==================== 图片管理 API ====================

// 上传图片
app.post('/api/images', requireAuth, upload.single('image'), (req, res) => {
    if (!req.file) {
        return res.status(400).json({ error: '没有上传文件' });
    }

    const { category, description } = req.body;
    const imageUrl = `/uploads/${req.file.filename}`;
});
