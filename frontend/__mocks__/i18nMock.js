const i18nMock = {
  use: () => i18nMock,
  init: () => i18nMock,
  t: (key) => key,
  changeLanguage: () => Promise.resolve(),
  language: 'zh',
  languages: ['zh', 'en'],
  services: {
    resourceStore: {
      data: {},
    },
  },
};

module.exports = i18nMock;