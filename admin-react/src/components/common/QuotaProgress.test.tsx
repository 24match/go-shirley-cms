/**
 * QuotaProgress 组件单元测试
 */

import React from 'react';
import { render, screen } from '@testing-library/react';
import QuotaProgress from './QuotaProgress';

describe('QuotaProgress', () => {
  describe('基础渲染', () => {
    it('应该渲染组件', () => {
      render(
        <QuotaProgress
          used={50}
          limit={100}
          resourceName="图片数量"
        />
      );

      expect(screen.getByText('图片数量')).toBeInTheDocument();
    });

    it('应该显示正确的使用量和限制', () => {
      render(
        <QuotaProgress
          used={50}
          limit={100}
          resourceName="存储空间"
          showLabels={true}
        />
      );

      expect(screen.getByText('50')).toBeInTheDocument();
      expect(screen.getByText('100')).toBeInTheDocument();
    });

    it('应该显示百分比', () => {
      render(
        <QuotaProgress
          used={50}
          limit={100}
          resourceName="测试资源"
        />
      );

      expect(screen.getByText('50%')).toBeInTheDocument();
    });
  });

  describe('状态显示', () => {
    it('正常状态应该显示绿色进度条', () => {
      const { container } = render(
        <QuotaProgress
          used={30}
          limit={100}
          resourceName="测试资源"
        />
      );

      const progressBar = container.querySelector('.quota-progress-bar');
      expect(progressBar).toHaveClass('quota-progress-normal');
    });

    it('警告状态应该显示黄色进度条（80% 以上）', () => {
      const { container } = render(
        <QuotaProgress
          used={85}
          limit={100}
          resourceName="测试资源"
        />
      );

      const progressBar = container.querySelector('.quota-progress-bar');
      expect(progressBar).toHaveClass('quota-progress-warning');
      
      // 应该显示警告图标
      expect(screen.getByText('⚡')).toBeInTheDocument();
    });

    it('错误状态应该显示红色进度条（100%）', () => {
      const { container } = render(
        <QuotaProgress
          used={100}
          limit={100}
          resourceName="测试资源"
        />
      );

      const progressBar = container.querySelector('.quota-progress-bar');
      expect(progressBar).toHaveClass('quota-progress-error');
      
      // 应该显示警告图标
      expect(screen.getByText('⚠️')).toBeInTheDocument();
    });

    it('无限制状态应该显示特殊样式', () => {
      const { container } = render(
        <QuotaProgress
          used={1000}
          limit={-1}
          resourceName="测试资源"
        />
      );

      const progressBar = container.querySelector('.quota-progress-bar');
      expect(progressBar).toHaveClass('quota-progress-unlimited');
      
      // 应该显示无限符号
      expect(screen.getByText('∞')).toBeInTheDocument();
      expect(screen.getByText('无限制')).toBeInTheDocument();
    });
  });

  describe('自定义阈值', () => {
    it('应该使用自定义警告阈值', () => {
      const { container } = render(
        <QuotaProgress
          used={60}
          limit={100}
          resourceName="测试资源"
          warningThreshold={50}
        />
      );

      const progressBar = container.querySelector('.quota-progress-bar');
      expect(progressBar).toHaveClass('quota-progress-warning');
    });

    it('应该使用自定义错误阈值', () => {
      const { container } = render(
        <QuotaProgress
          used={90}
          limit={100}
          resourceName="测试资源"
          errorThreshold={90}
        />
      );

      const progressBar = container.querySelector('.quota-progress-bar');
      expect(progressBar).toHaveClass('quota-progress-error');
    });
  });

  describe('数值格式化', () => {
    it('应该格式化大于 1000 的数值', () => {
      render(
        <QuotaProgress
          used={1500}
          limit={10000}
          resourceName="测试资源"
        />
      );

      expect(screen.getByText('1.5K')).toBeInTheDocument();
    });

    it('应该显示无限制符号', () => {
      render(
        <QuotaProgress
          used={5000}
          limit={-1}
          resourceName="测试资源"
        />
      );

      expect(screen.getByText('∞')).toBeInTheDocument();
    });
  });

  describe('隐藏标签', () => {
    it('当 showLabels 为 false 时不显示头部和底部', () => {
      const { container } = render(
        <QuotaProgress
          used={50}
          limit={100}
          resourceName="测试资源"
          showLabels={false}
        />
      );

      // 头部和底部不应该存在
      expect(container.querySelector('.quota-progress-header')).not.toBeInTheDocument();
      expect(container.querySelector('.quota-progress-footer')).not.toBeInTheDocument();
      
      // 进度条应该存在
      expect(container.querySelector('.quota-progress-bar')).toBeInTheDocument();
    });
  });

  describe('边界情况', () => {
    it('处理使用量为 0 的情况', () => {
      render(
        <QuotaProgress
          used={0}
          limit={100}
          resourceName="测试资源"
        />
      );

      expect(screen.getByText('0%')).toBeInTheDocument();
    });

    it('处理使用量超过限制的情况', () => {
      const { container } = render(
        <QuotaProgress
          used={150}
          limit={100}
          resourceName="测试资源"
        />
      );

      const progressBar = container.querySelector('.quota-progress-bar');
      expect(progressBar).toHaveClass('quota-progress-error');
      
      // 百分比应该限制在 100%
      expect(screen.getByText('100%')).toBeInTheDocument();
    });

    it('处理限制为 0 的情况', () => {
      render(
        <QuotaProgress
          used={0}
          limit={0}
          resourceName="测试资源"
        />
      );

      // 应该显示 100%（因为 0/0 会被处理）
      const percentageText = screen.getByText('100%');
      expect(percentageText).toBeInTheDocument();
    });
  });
});