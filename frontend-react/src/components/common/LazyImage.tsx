/**
 * LazyImage 组件
 * 支持懒加载的图片组件
 */

import React, { useEffect, useRef, useState } from 'react';

export interface LazyImageProps {
    dataSrc: string;
    alt: string;
    className?: string;
    placeholder?: string;
    style?: React.CSSProperties;
    onLoad?: () => void;
}

const LazyImage: React.FC<LazyImageProps> = ({
    dataSrc,
    alt,
    className = '',
    placeholder = 'data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7',
    style,
    onLoad
}) => {
    const imgRef = useRef<HTMLImageElement>(null);
    const [isLoaded, setIsLoaded] = useState(false);
    const [isInView, setIsInView] = useState(false);

    useEffect(() => {
        const observer = new IntersectionObserver(
            (entries) => {
                entries.forEach((entry) => {
                    if (entry.isIntersecting) {
                        setIsInView(true);
                        observer.disconnect();
                    }
                });
            },
            { rootMargin: '50px' }
        );

        if (imgRef.current) {
            observer.observe(imgRef.current);
        }

        return () => {
            observer.disconnect();
        };
    }, []);

    const handleLoad = () => {
        setIsLoaded(true);
        onLoad?.();
    };

    return (
        <img
            ref={imgRef}
            src={isInView ? dataSrc : placeholder}
            data-src={dataSrc}
            alt={alt}
            className={`lazy ${className} ${isLoaded ? 'loaded' : ''}`}
            style={style}
            onLoad={handleLoad}
        />
    );
};

export default LazyImage;