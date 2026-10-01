import React, { useEffect, useRef, useState } from 'react';
import './App.css';
import {
  QUADRATIC_EXAMPLE_CODE,
  TRIANGLE_EXAMPLE_CODE,
  TWENTY_ONE_EXAMPLE_CODE,
} from './legacyExamples';

interface RunResponse {
  sessionId: string;
}

interface StatusResponse {
  output?: string;
  needsInput?: boolean;
  completed?: boolean;
  error?: string;
}

const rawApiBaseUrl = process.env.REACT_APP_API_BASE_URL ?? 'http://127.0.0.1:5001';
const normalizedApiBaseUrl = /^https?:\/\//i.test(rawApiBaseUrl)
  ? rawApiBaseUrl
  : `http://${rawApiBaseUrl}`;
const API_BASE_URL = normalizedApiBaseUrl.replace(/\/+$/, '');

function apiUrl(path: string): string {
  return `${API_BASE_URL}${path}`;
}

function App() {
  const [code, setCode] = useState('');
  const [resultText, setResultText] = useState('執行結果將顯示在這裡');
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [isRunning, setIsRunning] = useState(false);
  const [showInputSection, setShowInputSection] = useState(false);
  const [userInput, setUserInput] = useState('');

  const inputRef = useRef<HTMLInputElement | null>(null);
  const pollTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const clearPollTimer = () => {
    if (pollTimeoutRef.current) {
      clearTimeout(pollTimeoutRef.current);
      pollTimeoutRef.current = null;
    }
  };

  const resetButton = () => {
    setIsRunning(false);
    setSessionId(null);
    clearPollTimer();
  };

  const pollStatus = async (sid: string) => {
    try {
      const response = await fetch(apiUrl(`/status/${sid}`));
      const data = (await response.json()) as StatusResponse;

      if (data.error) {
        setResultText(data.error);
        resetButton();
        return;
      }

      if (data.output) {
        setResultText((previous) => previous + data.output);
      }

      if (data.needsInput) {
        setShowInputSection(true);
      } else if (data.completed) {
        resetButton();
      } else {
        pollTimeoutRef.current = setTimeout(() => {
          void pollStatus(sid);
        }, 300);
      }
    } catch (_error) {
      setResultText('狀態輪詢失敗。');
      resetButton();
    }
  };

  const runCode = async () => {
    setIsRunning(true);
    setResultText('');
    setShowInputSection(false);

    try {
      const response = await fetch(apiUrl('/run'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ code }),
      });

      const resultData = (await response.json()) as RunResponse;
      setSessionId(resultData.sessionId);
      void pollStatus(resultData.sessionId);
    } catch (_error) {
      setResultText('前端請求錯誤或連線失敗。');
      resetButton();
    }
  };

  const stopExecution = async () => {
    if (!sessionId) {
      setResultText('目前沒有執行中的程式。');
      return;
    }

    try {
      const response = await fetch(apiUrl('/stop'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ sessionId }),
      });

      if (response.ok) {
        setResultText((previous) => `${previous}\n程式已強制終止。`);
        resetButton();
      } else {
        setResultText((previous) => `${previous}\n終止程式失敗。`);
      }
    } catch (error) {
      const message = error instanceof Error ? error.message : String(error);
      setResultText((previous) => `${previous}\n終止請求失敗：${message}`);
    }
  };

  const submitInput = async () => {
    const inputValue = userInput.trim();

    if (!inputValue || !sessionId) {
      alert('請輸入內容或確保 sessionId 存在');
      return;
    }

    try {
      const response = await fetch(apiUrl('/submit_input'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ sessionId, input: inputValue }),
      });

      if (!response.ok) {
        throw new Error('提交輸入失敗');
      }

      setUserInput('');
      void pollStatus(sessionId);
    } catch (error) {
      const message = error instanceof Error ? error.message : String(error);
      setResultText((previous) => `${previous}\n輸入處理失敗：${message}`);
    }
  };

  useEffect(() => {
    if (showInputSection) {
      inputRef.current?.focus();
    }
  }, [showInputSection]);

  useEffect(() => {
    return () => {
      clearPollTimer();
    };
  }, []);

  return (
    <div className="container">
      <div className="intro-section">
        <h2>基本介紹</h2>
        <p>這裡是基於python開發的直譯器。</p>
        <p>如果有任何問題請看下方的對照表。</p>
      </div>

      <h2>ACcode線上直譯器</h2>
      <textarea
        id="codeEditor"
        placeholder="在此輸入程式碼..."
        value={code}
        onChange={(event) => {
          setCode(event.target.value);
        }}
      />
      <button id="runCodeButton" onClick={() => void runCode()} disabled={isRunning}>
        {isRunning ? '執行中...' : '執行程式碼 ▶'}
      </button>
      <button id="stopButton" onClick={() => void stopExecution()}>
        強制結束 ⏹
      </button>
      <button
        id="triangleExampleButton"
        onClick={() => {
          setCode(TRIANGLE_EXAMPLE_CODE);
          setResultText('已載入三角形範例程式，請點擊「執行程式碼」運行。');
        }}
      >
        載入三角形範例
      </button>
      <button
        id="quadraticExampleButton"
        onClick={() => {
          setCode(QUADRATIC_EXAMPLE_CODE);
          setResultText('已載入一元二次方程式範例程式，請點擊「執行程式碼」運行。');
        }}
      >
        載入二次方程範例
      </button>
      <button
        id="twentyOneExampleButton"
        onClick={() => {
          setCode(TWENTY_ONE_EXAMPLE_CODE);
          setResultText('已載入搶21點範例程式，請點擊「執行程式碼」運行。');
        }}
      >
        載入搶21點範例
      </button>

      <pre id="executionResult">{resultText}</pre>
      <div id="inputSection" style={{ display: showInputSection ? 'block' : 'none' }}>
        <input
          ref={inputRef}
          type="text"
          id="userInput"
          placeholder="請輸入..."
          value={userInput}
          onChange={(event) => {
            setUserInput(event.target.value);
          }}
        />
        <button id="submitInput" onClick={() => void submitInput()}>
          送出輸入
        </button>
      </div>

      <table className="syntax-table">
        <thead>
          <tr>
            <th>Python語法</th>
            <th>ACcode語法</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>print(...)</td>
            <td>要求在終端機（或控制台）輸出 ... 。</td>
          </tr>
          <tr>
            <td>... = input()</td>
            <td>讀取使用者輸入字串，並儲存至變數「...」</td>
          </tr>
          <tr>
            <td>... = ....</td>
            <td>設定一指定名稱變數的值，名稱與值分別為 ... ， ....</td>
          </tr>
          <tr>
            <td>... = [None for _ in range(....)]</td>
            <td>初始化一個指定名稱和長度的陣列，名稱和長度分別為 ... ， ....</td>
          </tr>
          <tr>
            <td>if ... :</td>
            <td>如果右方判斷式結果為真，執行以下特定操作</td>
          </tr>
          <tr>
            <td>while ... :</td>
            <td>重複並持續驗證右方表達式之真實性，條件滿足時持續執行特定操作 ...</td>
          </tr>
          <tr>
            <td>
              <i>取消縮排</i>
            </td>
            <td>結束以上判斷式或迴圈</td>
          </tr>
          <tr>
            <td>break</td>
            <td>無視迴圈判斷式要求，直接跳脫迴圈</td>
          </tr>
          <tr>
            <td>continue</td>
            <td>捨棄以下迴圈內容，直接繼續下一輪迴圈</td>
          </tr>
          <tr>
            <td>exit()</td>
            <td>無視所有指令，直接退出程式</td>
          </tr>
          <tr>
            <td>+, -, *, /, %</td>
            <td>加上, 減去, 乘以, 除以, 取模</td>
          </tr>
          <tr>
            <td>=, &gt;, &lt;</td>
            <td>等於, 大於, 小於</td>
          </tr>
          <tr>
            <td>float(...), str(...)</td>
            <td>... 轉型 「數字」, ... 轉型 「字串」</td>
          </tr>
          <tr>
            <td>
              ... <i>(變數)</i>, ...[....]
              <i>(清單)</i>
            </td>
            <td>變數「...」, 陣列「...」的索引（....）</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
}

export default App;
