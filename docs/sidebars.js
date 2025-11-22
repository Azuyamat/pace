/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  docsSidebar: [
    'intro',
    {
      type: 'category',
      label: 'Getting Started',
      items: ['installation', 'quick-start'],
    },
    {
      type: 'category',
      label: 'Commands',
      items: ['commands/add', 'commands/list', 'commands/complete', 'commands/delete'],
    },
  ],
};

module.exports = sidebars;
