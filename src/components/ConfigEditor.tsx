import React, { ChangeEvent } from 'react';
import { InlineField, Input, SecretInput } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { UpdownDataSourceOptions, UpdownSecureJsonData } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<UpdownDataSourceOptions, UpdownSecureJsonData> {}

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;
  const { jsonData, secureJsonFields, secureJsonData } = options;

  const onApiUrlChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        apiUrl: event.target.value,
      },
    });
  };

  // Secure field (only sent to the backend)
  const onAPIKeyChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        apiKey: event.target.value,
      },
    });
  };

  const onResetAPIKey = () => {
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        apiKey: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        apiKey: '',
      },
    });
  };

  return (
    <>
      <InlineField 
        label="API URL" 
        labelWidth={14} 
        interactive 
        tooltip={'Base URL for the UpDown.io API (optional, default: https://updown.io/api)'}
      >
        <Input
          id="config-editor-api-url"
          onChange={onApiUrlChange}
          value={jsonData.apiUrl || 'https://updown.io/api'}
          placeholder="https://updown.io/api"
          width={40}
        />
      </InlineField>
      <InlineField 
        label="API Key" 
        labelWidth={14} 
        interactive 
        tooltip={'Your UpDown.io API key (available in your UpDown.io account)'}
      >
        <SecretInput
          required
          id="config-editor-api-key"
          isConfigured={secureJsonFields.apiKey}
          value={secureJsonData?.apiKey}
          placeholder="Enter your UpDown.io API key"
          width={40}
          onReset={onResetAPIKey}
          onChange={onAPIKeyChange}
        />
      </InlineField>
    </>
  );
}
