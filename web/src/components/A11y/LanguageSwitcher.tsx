import React, { useState } from 'react';
import { Button, Dropdown, Space } from 'antd';
import { GlobalOutlined, DownOutlined } from '@ant-design/icons';
import type { MenuProps } from 'antd';
import { useI18n } from '@/hooks/useI18n';

export const LanguageSwitcher: React.FC = () => {
  const { t, language, changeLanguage, languages } = useI18n();
  const [loading, setLoading] = useState(false);

  const currentLanguage = languages.find(lang => lang.code === language);

  const handleLanguageChange = async (langCode: string) => {
    if (langCode === language) return;
    
    setLoading(true);
    try {
      await changeLanguage(langCode);
    } finally {
      setLoading(false);
    }
  };

  const menuItems: MenuProps['items'] = languages.map(lang => ({
    key: lang.code,
    label: (
      <div className="flex items-center justify-between min-w-[120px]">
        <span>{lang.nativeName}</span>
        {lang.code === language && (
          <span className="text-blue-500 ml-2">✓</span>
        )}
      </div>
    ),
    onClick: () => handleLanguageChange(lang.code),
  }));

  return (
    <Dropdown
      menu={{ items: menuItems }}
      trigger={['click']}
      placement="bottomRight"
      arrow
    >
      <Button
        type="text"
        loading={loading}
        aria-label={t('accessibility.languageSwitcher')}
        aria-expanded="false"
        aria-haspopup="true"
        className="flex items-center"
      >
        <Space>
          <GlobalOutlined />
          <span className="hidden sm:inline">
            {currentLanguage?.nativeName}
          </span>
          <DownOutlined className="text-xs" />
        </Space>
      </Button>
    </Dropdown>
  );
};