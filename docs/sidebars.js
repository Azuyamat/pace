/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  docsSidebar: [
    'intro',
    {
      type: 'category',
      label: 'Getting Started',
      items: ['installation', 'quick-start'],
    },
    'configuration',
    {
      type: 'category',
      label: 'Commands',
      items: ['commands/run', 'commands/watch', 'commands/list', 'commands/update', 'commands/version'],
    },
    'examples',
  ],
};

module.exports = sidebars;
