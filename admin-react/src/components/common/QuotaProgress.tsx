/**
 * 配额使用进度条组件
 * 用于展示资源使用情况，支持警告和错误状态
 */

import React from 'react';

interface QuotaProgressProps {
  /** 已使用量 */
  used: number;
  /** 配额限制 */
  limit: number;
  /** 资源名称 */
  resourceName: string;
  /** 是否显示数值标签 */
  showLabels?: boolean;
  /** 自定义警告阈值（默认 80%） */
  warningThreshold?: number;
  /** 自定义错误阈值（默认 100%） */
  errorThreshold?: number;
}

const QuotaProgress: React.FC<QuotaProgressProps> = ({
  used,
  limit,
  resourceName,
  showLabels = true,
  warningThreshold = 80,
  errorThreshold = 100,
}) => {
  // 计算使用百分比
  const isUnlimited = limit === -1;
  const percentage = isUnlimited ? 0 : Math.min((used / limit) * 100, 100);
  
  // 确定状态
  const getStatus = () => {
    if (isUnlimited) return 'unlimited';
    if (percentage >= errorThreshold) return 'error';
    if (percentage >= warningThreshold) return 'warning';
    return 'normal';
  };

  const status = getStatus();

  // 获取状态对应的样式类
  const getStatusClass = () => {
    switch (status) {
      case 'error':
        return 'quota-progress-error';
      case 'warning':
        return 'quota-progress-warning';
      case 'unlimited':
        return 'quota-progress-unlimited';
      default:
        return 'quota-progress-normal';
    }
  };

  // 获取状态对应的图标
  const getStatusIcon = () => {
    switch (status) {
      case 'error':
        return '⚠️';
      case 'warning':
        return '⚡';
      case 'unlimited':
        return '∞';
      default:
        return '✓';
    }
  };

  // 格式化数值显示
  const formatValue = (value: number) => {
    if (value === -1) return '无限制';
    if (value >= 1000) return `${(value / 1000).toFixed(1)}K`;
    return value.toString();
  };

  return (
    <div className="quota-progress-container">
      {showLabels && (
        <div className="quota-progress-header">
          <span className="quota-progress-name">{resourceName}</span>
          <span className={`quota-progress-status ${status}`}>
            {getStatusIcon()} {isUnlimited ? '无限制' : `${Math.round(percentage)}%`}
          </span>
        </div>
      )}
      
      <div className={`quota-progress-bar ${getStatusClass()}`}>
        <div 
          className="quota-progress-fill"
          style={{ 
            width: isUnlimited ? '100%' : `${percentage}%`,
            background: isUnlimited 
              ? 'linear-gradient(90deg, #10b981 0%, #34d399 100%)'
              : status === 'error'
                ? 'linear-gradient(90deg, #ef4444 0%, #f87171 100%)'
                : status === 'warning'
                  ? 'linear-gradient(90deg, #f59e0b 0%, #fbbf24 100%)'
                  : 'linear-gradient(90deg, #10b981 0%, #34d399 100%)'
          }}
        />
      </div>
      
      {showLabels && (
        <div className="quota-progress-footer">
          <span className="quota-progress-used">{formatValue(used)}</span>
          <span className="quota-progress-limit">
            {isUnlimited ? '∞' : `/ ${formatValue(limit)}`}
          </span>
        </div>
      )}
    </div>
  );
};

export default QuotaProgress;