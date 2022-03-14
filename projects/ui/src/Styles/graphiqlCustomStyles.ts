import { css } from '@emotion/core';
import { colors } from './colors';

const graphiqlCustomStyles = css`
  .graphiql-container {
    button {
      font-family: 'Proxima Nova', 'Open Sans', 'Helvetica', 'Arial',
        'sans-serif';
    }

    // ========================== //
    // TITLE BAR
    .editorWrap .topBarWrap .topBar {
      background: #fff;
      height: max-content;
      padding: 8px 0 9px 10px;
      .title {
        opacity: 0.5;
        transform: scale(0.9);
      }
    }

    // ========================== //
    // RESULT WINDOW
    .result-window .CodeMirror {
      background: white;
    }

    // ========================== //
    // TOOLBAR + DOC EXPLORER
    button.docExplorerShow,
    button.doc-explorer-back,
    .editorWrap .topBarWrap .topBar .toolbar button.toolbar-button {
      padding: 5px 10px 5px;
      box-shadow: none;
      background: none;
      color: ${colors.seaBlue};
      &:hover {
        color: ${colors.lakeBlue};
      }
      &:active {
        color: ${colors.pondBlue};
      }
    }
    .editorWrap .topBarWrap .topBar .toolbar button.toolbar-button {
      border: none;
    }
    button.docExplorerShow {
      border-left: none;
    }
    button.doc-explorer-back {
      margin-left: 0px;
    }
    button.docExplorerHide {
      padding: 10px;
    }
    button.docExplorerShow:before,
    button.doc-explorer-back:before {
      margin-bottom: 1px;
      border-color: ${colors.seaBlue};
    }
    .docExplorerWrap {
      box-shadow: none;
    }
    .doc-explorer-title-bar {
      display: flex;
      align-items: center;
      justify-content: center;
      padding-top: 0px;
      padding-bottom: 0px;
    }

    // ========================== //
    // GUTTERS/RESIZERS
    .docExplorerWrap .docExplorerResizer {
      border-left: 1px solid #e0e0e0;
      border-right: 1px solid #e0e0e0;
      width: 0.7em;
      padding-left: 3px;
      margin-left: -3px;
      box-sizing: content-box;
    }
    .docExplorerWrap .docExplorerResizer,
    .CodeMirror-gutter.CodeMirror-foldgutter:only-child {
      display: flex;
      align-items: center;
      justify-content: center;
      &:after {
        content: '';
        height: 60px;
        max-height: 40%;
        width: 4px;
        margin-right: 2px;
        border-radius: 100px;
        transition: 0.2s background-color;
        background-color: transparent;
      }
      &:active:after,
      &:hover:after {
        background-color: ${colors.marchGrey};
        background-color: ${colors.mayGrey};
      }
    }
    .docExplorerWrap .docExplorerResizer,
    .CodeMirror-gutter.CodeMirror-foldgutter,
    .CodeMirror-gutter.CodeMirror-linenumbers {
      background-color: white;
    }
    .secondary-editor-title.variable-editor-title {
      background-color: ${colors.januaryGrey};
    }

    // ========================== //
    // EXECUTE/RUN QUERY BUTTON //
    .execute-button-wrap {
      height: unset;
      order: 2;
      margin: 0 0 0 1rem;
      button {
        // Content ---
        transform: scale(0.95);
        svg {
          transform: translateX(-10px) scale(0.8);
          order: -1;
        }
        &:before {
          content: 'Run';
          transform: translateX(-10px);
        }
        // Sizing ---
        width: 70px;
        height: 30px;
        padding-left: 0.8rem;
        display: flex;
        align-items: center;
        // Colors ---
        justify-content: center;
        fill: ${colors.darkOceanBlue};
        stroke: ${colors.darkOceanBlue};
        color: ${colors.darkOceanBlue};
        background: ${colors.splashBlue};
        border: 1px solid ${colors.oceanBlue};
        border-radius: 20px;
        box-shadow: none;
        // Active + Hover ---
        transition: 0.1s transform;
        &:hover {
          transform: none;
          background: ${colors.pondBlue};
        }
        &:active {
          background: ${colors.lakeBlue};
        }
      }
    }
  }
`;

export default graphiqlCustomStyles;
