import React, { useState, useEffect } from "react";

const Contents = () => {
  const [messages, setMessages] = useState([]);
  const [now, setNow] = useState("");
  const [send, setSend] = useState("");
  const channels = ["channel1", "channel2", "channel3", "channel4"];

  return (
    <>
      <main>
        <ul>
            {messages.map((message) => (
                <li>{message.content}{!(message.posted == message.fixed) ? <>(編集済み)</> : null}</li>
            ))}
        </ul>
        </main>
        <aside>
            {now ? <p>あなたの現在のチャンネル：{now}</p> : <p>チャンネルを選択してください。</p>}
            {channels.map((channel) => (
                <button onClick={() => setNow(channel)}>{channel}</button>
            ))}
        </aside>
        <footer>
            {now ? <><textarea
                onChange={(e) => setSend(e.target.value)}
                ></textarea>
                <button>送信</button>
                </> : null
            }
        </footer>
    </>
  );
};

export default Contents;