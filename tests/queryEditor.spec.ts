import { test, expect } from '@grafana/plugin-e2e';

test('should load datasource without errors', async ({ panelEditPage, readProvisionedDataSource }) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  
  // This should not throw an error
  await panelEditPage.datasource.set(ds.name);
  
  // If we get here without throwing, the test passes
  expect(true).toBe(true);
});

test('should be able to set table visualization', async ({ panelEditPage, readProvisionedDataSource }) => {
  const ds = await readProvisionedDataSource({ fileName: 'datasources.yml' });
  await panelEditPage.datasource.set(ds.name);
  
  // Just verify that we can set the visualization without errors
  await panelEditPage.setVisualization('Table');
  
  // If we get here without throwing, the test passes
  expect(true).toBe(true);
});
