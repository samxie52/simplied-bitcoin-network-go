import { useTranslation } from 'react-i18next';
import { useCallback } from 'react';

export interface UseI18nReturn {
  t: (key: string, options?: any) => string;
  language: string;
  changeLanguage: (lng: string) => Promise<void>;
  languages: Array<{ code: string; name: string; nativeName: string }>;
  isRTL: boolean;
}

const SUPPORTED_LANGUAGES = [
  { code: 'zh-CN', name: 'Chinese', nativeName: '中文' },
  { code: 'en-US', name: 'English', nativeName: 'English' },
];

export const useI18n = (): UseI18nReturn => {
  const { t, i18n } = useTranslation();

  const changeLanguage = useCallback(async (lng: string) => {
    try {
      await i18n.changeLanguage(lng);
      // 更新HTML lang属性
      document.documentElement.lang = lng;
      // 更新页面方向
      document.documentElement.dir = lng === 'ar' ? 'rtl' : 'ltr';
    } catch (error) {
      console.error('Failed to change language:', error);
    }
  }, [i18n]);

  return {
    t,
    language: i18n.language,
    changeLanguage,
    languages: SUPPORTED_LANGUAGES,
    isRTL: i18n.dir() === 'rtl',
  };
};