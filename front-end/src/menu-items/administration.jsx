// assets
import { IconUsers, IconSignatureOff, IconPaperclip } from '@tabler/icons-react';

const administration = {
  id: 'administration',
  title: 'Administración',
  type: 'group',
  children: [
    {
      id: 'users',
      title: 'Usuarios',
      type: 'item',
      url: '/users',
      icon: IconUsers,
      breadcrumbs: false
    },
    {
      id: 'signatures',
      title: 'Firmas',
      type: 'item',
      url: '/signatures',
      icon: IconSignatureOff,
      breadcrumbs: false
    },
    {
      id: 'documents',
      title: 'Documentos Capturados',
      type: 'item',
      url: '/documents',
      icon: IconPaperclip,
      breadcrumbs: false
    }
  ]
};

export default administration;
