// assets
import { IconSignature, IconScan } from '@tabler/icons-react';

const signature = {
  id: 'sample-docs-roadmap',
  title: 'Firmas',
  type: 'group',
  children: [
    {
      id: 'capture',
      title: 'Escanear documento',
      type: 'item',
      url: '/scan',
      icon: IconScan,
      breadcrumbs: false
    },
    {
      id: 'signature',
      title: 'Mi firma',
      type: 'item',
      url: '/signature',
      icon: IconSignature,
      breadcrumbs: false
    },
  ]
};

export default signature;
