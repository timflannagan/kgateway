import * as React from 'react';
import styled from '@emotion/styled/macro';
import AceEditor, { IAceEditorProps } from 'react-ace';
/*
  These imports are needed for syntax highlighting and snippets. DO NOT REMOVE.
*/
import 'ace-builds/src-noconflict/ext-language_tools';
import 'ace-builds/src-noconflict/ext-searchbox';
import 'ace-builds/src-noconflict/mode-yaml';
import 'ace-builds/src-noconflict/mode-html';
import 'ace-builds/src-noconflict/mode-graphqlschema';
import 'ace-builds/src-noconflict/snippets/yaml';
import 'ace-builds/src-noconflict/snippets/graphqlschema';
import 'ace-builds/src-noconflict/theme-chrome';
import 'ace-builds/webpack-resolver';
import { colors } from 'Styles/colors';
import { SoloDropdown } from './SoloDropdown';

export const Label = styled.label`
  display: block;
  color: ${colors.novemberGrey};
  font-size: 16px;
  margin-bottom: 10px;
  font-weight: 500;
`;

const StyledAceEditor = styled(AceEditor)`
  .ace_editor span,
  .ace_editor textarea {
    font-size: 16px;
    font-family: 'monospace';
  }
`;

export interface SoloFormVisualEditorProps extends IAceEditorProps {
  name: string; // the name of this field in Formik
  title?: string; // display name of the field
}

const VisualEditor = (props: SoloFormVisualEditorProps) => {
  const { name, title, value, ...rest } = props;

  // Commenting this out since it looks strange when multiple editors are on one page.
  // TODO: Make a context API that stores lightweight global settings like this.
  // const lsKey = 'ace-keyboard-handler';
  // const lsKeyboardHander = localStorage.getItem(lsKey);
  // const [keyboardHandler, setKeyboardHandler] = React.useState(
  //   lsKeyboardHander ? lsKeyboardHander : 'keyboard'
  // );

  return (
    <div className='relative'>
      {title && <Label>{title}</Label>}

      <StyledAceEditor
        // keyboardHandler={keyboardHandler}
        mode={rest.mode ?? 'yaml'}
        theme='chrome'
        name={name ?? title}
        style={{
          maxWidth: '40vw',
          maxHeight: '25vh',
          // cursor: 'text',
        }}
        onChange={rest.onChange}
        focus={true}
        onInput={rest.onInput}
        fontSize={14}
        showPrintMargin={false}
        showGutter={true}
        highlightActiveLine={true}
        value={value}
        readOnly={false}
        setOptions={{
          highlightGutterLine: true,
          showGutter: true,
          cursorStyle: 'wide',
          fontFamily: 'monospace',
          enableBasicAutocompletion: true,
          enableLiveAutocompletion: true,
          showLineNumbers: true,
          tabSize: 2,
        }}
        {...rest}
      />

      {/* <div
        style={{
          position: 'absolute',
          right: '0px',
          bottom: '0px',
        }}>
        <SoloDropdown
          value={keyboardHandler}
          options={[
            { key: 'keyboard', value: 'keyboard' },
            { key: 'vim', value: 'vim' },
            { key: 'emacs', value: 'emacs' },
          ]}
          onChange={e => {
            const newKeyboardHandler = e as string;
            setKeyboardHandler(newKeyboardHandler);
            localStorage.setItem(lsKey, newKeyboardHandler);
          }}
        />
      </div> */}
    </div>
  );
};

export { VisualEditor as default };
