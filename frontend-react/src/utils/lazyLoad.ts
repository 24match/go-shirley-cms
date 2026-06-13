/**
 * 图片懒加载工具
 * 使用 IntersectionObserver 实现图片懒加载
 */

export interface LazyLoadOptions {
    rootMargin?: string;
    threshold?: number;
    placeholder?: string;
}

const defaultOptions: LazyLoadOptions = {
    rootMargin: '50px',
    threshold: 0,
    placeholder: 'data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7'
};

/**
 * 懒加载图片
 * @param img 图片元素
 * @param options 配置选项
 */
export function lazyLoadImage(
    img: HTMLImageElement, 
    options: LazyLoadOptions = defaultOptions
): void {
    const dataSrc = img.dataset.src;
    if (!dataSrc) return;

    const observer = new IntersectionObserver((entries, observer) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                const img = entry.target as HTMLImageElement;
                img.src = img.dataset.src!;
                img.classList.remove('lazy');
                observer.unobserve(img);
            }
        });
    }, {
        rootMargin: options.rootMargin,
        threshold: options.threshold
    });

    observer.observe(img);
}

/**
 * 懒加载多个图片
 */
export function lazyLoadImages(
    selector: string = 'img[data-src]',
    options: LazyLoadOptions = defaultOptions
): void {
    const images = document.querySelectorAll<HTMLImageElement>(selector);
    images.forEach(img => lazyLoadImage(img, options));
}

/**
 * 懒加载 Hook 辅助函数
 * 用于在 React 组件中初始化懒加载
 */
export function initLazyLoading(
    container: Element | Document = document,
    options: LazyLoadOptions = defaultOptions
): () => void {
    const images = container.querySelectorAll<HTMLImageElement>('img[data-src]');
    images.forEach(img => lazyLoadImage(img, options));

    // 返回清理函数
    return () => {
        // 可以在这里添加清理逻辑
    };
}