import { IconChartBar } from '@tabler/icons-react';

const dashboard = {
  id: 'dashboard',
  title: 'Dashboard',
  type: 'group',
  children: [
    {
      id: 'default',
      title: 'Dashboard',
      type: 'item',
      url: '/dashboard',
      icon: IconChartBar,
      breadcrumbs: false
    }
  ]
};

export default dashboard;
