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
        tooltip={'URL de base de l\'API UpDown.io (optionnel, par défaut: https://updown.io/api)'}
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
        label="Clé API" 
        labelWidth={14} 
        interactive 
        tooltip={'Votre clé API UpDown.io (disponible dans votre compte UpDown.io)'}
      >
        <SecretInput
          required
          id="config-editor-api-key"
          isConfigured={secureJsonFields.apiKey}
          value={secureJsonData?.apiKey}
          placeholder="Entrez votre clé API UpDown.io"
          width={40}
          onReset={onResetAPIKey}
          onChange={onAPIKeyChange}
        />
      </InlineField>
    </>
  );
}
