import React, { useCallback, useEffect, useState } from 'react';
import { InlineField, Select, Stack } from '@grafana/ui';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { DataSource } from '../datasource';
import { UpdownDataSourceOptions, UpdownQuery, UpdownCheck } from '../types';

type Props = QueryEditorProps<DataSource, UpdownQuery, UpdownDataSourceOptions>;

const queryTypeOptions: Array<SelectableValue<string>> = [
  { label: 'Service List', value: 'checks', description: 'Retrieve the list of all your monitored services' },
  { label: 'Performance Overview', value: 'performance', description: 'Complete performance dashboard with uptime, response time and status' },
  { label: 'Service Metrics', value: 'metrics', description: 'Retrieve detailed metrics for a service' },
  { label: 'Downtime', value: 'downtimes', description: 'Retrieve downtime history for a service' },
  { label: 'Uptime', value: 'uptime', description: 'Retrieve uptime percentage for a service' },
];

const groupByOptions: Array<SelectableValue<string>> = [
  { label: 'By time', value: 'time', description: 'Group data by time intervals' },
  { label: 'By location', value: 'host', description: 'Group data by monitoring server' },
];

export function QueryEditor({ query, onChange, onRunQuery, datasource }: Props) {
  const [checks, setChecks] = useState<UpdownCheck[]>([]);
  const [loading, setLoading] = useState(false);

  const loadChecks = useCallback(async () => {
    try {
      setLoading(true);
      const result = await datasource.getChecks();
      setChecks(result);
    } catch (error) {
      console.error('Error loading checks:', error);
    } finally {
      setLoading(false);
    }
  }, [datasource]);

  // Load checks list when component mounts
  useEffect(() => {
    loadChecks();
  }, [loadChecks]);

  const onQueryTypeChange = (selection: SelectableValue<string>) => {
    onChange({ 
      ...query, 
      queryType: selection.value || 'checks',
      // Reset checkToken if changing to 'checks'
      checkToken: selection.value === 'checks' ? undefined : query.checkToken
    });
    onRunQuery();
  };

  const onCheckTokenChange = (selection: SelectableValue<string>) => {
    onChange({ ...query, checkToken: selection.value });
    onRunQuery();
  };

  const onGroupByChange = (selection: SelectableValue<string>) => {
    onChange({ ...query, groupBy: selection.value });
    onRunQuery();
  };

  const { queryType, checkToken, groupBy } = query;

  // Options for check selection
  const checkOptions: Array<SelectableValue<string>> = checks.map(check => ({
    label: check.alias || check.url,
    value: check.token,
    description: `${check.url} (${check.uptime}% uptime)`,
  }));

  const needsCheckToken = queryType && ['metrics', 'downtimes', 'uptime'].includes(queryType);
  const needsGroupBy = queryType === 'metrics';

  return (
    <Stack gap={0}>
      <InlineField label="Data type" labelWidth={16}>
        <Select
          width={40}
          value={queryType}
          placeholder="Select data type"
          options={queryTypeOptions}
          onChange={onQueryTypeChange}
        />
      </InlineField>

      {needsCheckToken && (
        <InlineField label="Service" labelWidth={16}>
          <Select
            width={40}
            value={checkToken}
            placeholder={loading ? "Loading..." : "Select a service"}
            options={checkOptions}
            onChange={onCheckTokenChange}
            isLoading={loading}
          />
        </InlineField>
      )}

      {needsGroupBy && (
        <InlineField label="Grouping" labelWidth={16} tooltip="How to group the data">
          <Select
            width={40}
            value={groupBy}
            placeholder="Group by..."
            options={groupByOptions}
            onChange={onGroupByChange}
          />
        </InlineField>
      )}
    </Stack>
  );
}
