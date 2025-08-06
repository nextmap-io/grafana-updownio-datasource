import { test, expect } from '@grafana/plugin-e2e';

test('smoke: should render query editor', async ({ panelEditPage, readProvisionedDataSource }) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await panelEditPage.datasource.set(ds.name);
  await expect(panelEditPage.getQueryEditorRow('A').getByRole('combobox')).toBeVisible();
});

test('should show default query type selection', async ({
  panelEditPage,
  readProvisionedDataSource,
}) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await panelEditPage.datasource.set(ds.name);
  
  // Should show the Select data type dropdown
  await expect(panelEditPage.getQueryEditorRow('A').getByText('Select data type')).toBeVisible();
});

test('should render query editor interface', async ({ panelEditPage, readProvisionedDataSource }) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await panelEditPage.datasource.set(ds.name);
  await panelEditPage.setVisualization('Table');
  
  // Just verify the interface loads without trying to execute queries
  await expect(panelEditPage.getQueryEditorRow('A').getByRole('combobox')).toBeVisible();
});
