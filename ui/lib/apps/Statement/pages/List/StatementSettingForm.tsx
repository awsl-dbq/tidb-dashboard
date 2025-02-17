import React, { useState, useCallback } from 'react'
import {
  Form,
  Skeleton,
  Switch,
  Input,
  Slider,
  Space,
  Button,
  Modal,
} from 'antd'
import { ExclamationCircleOutlined } from '@ant-design/icons'
import { useTranslation } from 'react-i18next'
import client, { StatementEditableConfig } from '@lib/client'
import { useClientRequest } from '@lib/utils/useClientRequest'
import { ErrorBar } from '@lib/components'
import { useIsWriteable } from '@lib/utils/store'

interface Props {
  onClose: () => void
  onConfigUpdated: () => any
}

const convertArrToObj = (arr: number[]) =>
  arr.reduce((acc, cur) => {
    acc[cur] = cur
    return acc
  }, {})

function StatementSettingForm({ onClose, onConfigUpdated }: Props) {
  const [submitting, setSubmitting] = useState(false)
  const { t } = useTranslation()
  const isWriteable = useIsWriteable()

  const {
    data: initialConfig,
    isLoading: loading,
    error,
  } = useClientRequest((reqConfig) =>
    client.getInstance().statementsConfigGet(reqConfig)
  )

  const handleSubmit = useCallback(
    (values) => {
      async function updateConfig(values) {
        const newConfig: StatementEditableConfig = {
          enable: values.enable,
          max_size: values.max_size,
          refresh_interval: values.refresh_interval * 60,
          history_size: values.history_size,
          internal_query: values.internal_query,
        }
        try {
          setSubmitting(true)
          await client.getInstance().statementsConfigPost(newConfig)
          onClose()
          onConfigUpdated()
        } finally {
          setSubmitting(false)
        }
      }

      if (!values.enable) {
        // warning
        Modal.confirm({
          title: t('statement.settings.close_statement'),
          icon: <ExclamationCircleOutlined />,
          content: t('statement.settings.close_statement_warning'),
          okText: t('statement.settings.actions.close'),
          cancelText: t('statement.settings.actions.cancel'),
          okButtonProps: { danger: true },
          onOk: () => updateConfig(values),
        })
      } else {
        updateConfig(values)
      }
    },
    [t, onClose, onConfigUpdated]
  )

  return (
    <>
      {error && <ErrorBar errors={[error]} />}
      {loading && <Skeleton active={true} paragraph={{ rows: 5 }} />}
      {!loading && initialConfig && (
        <Form
          layout="vertical"
          initialValues={{
            ...initialConfig,
            refresh_interval: Math.floor(
              (initialConfig.refresh_interval ?? 0) / 60
            ),
          }}
          onFinish={handleSubmit}
        >
          <Form.Item
            valuePropName="checked"
            label={t('statement.settings.switch')}
            extra={t('statement.settings.switch_tooltip')}
          >
            <Form.Item noStyle name="enable" valuePropName="checked">
              <Switch disabled={!isWriteable} />
            </Form.Item>
          </Form.Item>
          <Form.Item
            noStyle
            shouldUpdate={(prev, cur) => prev.enable !== cur.enable}
          >
            {({ getFieldValue }) =>
              getFieldValue('enable') && (
                <>
                  <Form.Item
                    label={t('statement.settings.max_size')}
                    extra={t('statement.settings.max_size_tooltip')}
                  >
                    <Input.Group>
                      <Form.Item noStyle name="max_size">
                        <Slider
                          disabled={!isWriteable}
                          min={200}
                          max={5000}
                          step={100}
                          marks={convertArrToObj([200, 1000, 2000, 5000])}
                        />
                      </Form.Item>
                    </Input.Group>
                  </Form.Item>
                  <Form.Item
                    label={t('statement.settings.refresh_interval')}
                    extra={t('statement.settings.refresh_interval_tooltip')}
                  >
                    <Input.Group>
                      <Form.Item noStyle name="refresh_interval">
                        <Slider
                          disabled={!isWriteable}
                          min={1}
                          max={60}
                          step={null}
                          marks={convertArrToObj([1, 5, 15, 30, 60])}
                        />
                      </Form.Item>
                    </Input.Group>
                  </Form.Item>
                  <Form.Item
                    label={t('statement.settings.history_size')}
                    extra={t('statement.settings.history_size_tooltip')}
                  >
                    <Input.Group>
                      <Form.Item noStyle name="history_size">
                        <Slider
                          disabled={!isWriteable}
                          min={1}
                          max={255}
                          marks={convertArrToObj([1, 255])}
                        />
                      </Form.Item>
                    </Input.Group>
                  </Form.Item>
                  <Form.Item
                    label={t('statement.settings.keep_duration')}
                    extra={t('statement.settings.keep_duration_tooltip')}
                    shouldUpdate={(prev, cur) =>
                      prev.refresh_interval !== cur.refresh_interval ||
                      prev.history_size !== cur.history_size
                    }
                  >
                    {({ getFieldValue }) => {
                      const refreshInterval =
                        getFieldValue('refresh_interval') || 0
                      const historySize = getFieldValue('history_size') || 0
                      const totalMins = refreshInterval * historySize
                      const day = Math.floor(totalMins / (24 * 60))
                      const hour = Math.floor((totalMins - day * 24 * 60) / 60)
                      const min = totalMins - day * 24 * 60 - hour * 60
                      return `${day} day ${hour} hour ${min} min`
                    }}
                  </Form.Item>
                  <Form.Item
                    label={t('statement.settings.internal_query')}
                    extra={t('statement.settings.internal_query_tooltip')}
                    name="internal_query"
                    valuePropName="checked"
                  >
                    <Switch disabled={!isWriteable} />
                  </Form.Item>
                </>
              )
            }
          </Form.Item>
          <Form.Item>
            <Space>
              <Button
                type="primary"
                htmlType="submit"
                loading={submitting}
                disabled={!isWriteable}
              >
                {t('statement.settings.actions.save')}
              </Button>
              <Button onClick={onClose}>
                {t('statement.settings.actions.cancel')}
              </Button>
            </Space>
          </Form.Item>
        </Form>
      )}
    </>
  )
}

export default StatementSettingForm
