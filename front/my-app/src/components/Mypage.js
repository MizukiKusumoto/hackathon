import React, { useState, useEffect } from "react";
import { onAuthStateChanged, signOut } from "firebase/auth";
import { auth } from "../firebase.js";
import { useNavigate, Navigate } from "react-router-dom";

const Mypage = () => {
  const [user, setUser] = useState("");
  const [loading, setLoading] = useState(true);
  const [messages, setMessages] = useState([]);
  const [now, setNow] = useState("");
  const [send, setSend] = useState("");
  const channels = ["channel1", "channel2", "channel3", "channel4"];

  useEffect(() => {
    onAuthStateChanged(auth, (currentUser) => {
      setUser(currentUser);
      setLoading(false);
    });
  }, []);

  const navigate = useNavigate();

  const logout = async () => {
    await signOut(auth);
    navigate("/login/");
  }

  return (
    <>
      {!loading && (
        <>
          {!user ? (
            <Navigate to={`/login/`} />
          ) : (
            <>
              <header>
                <h1 onClick={() => setMessages([{content: "あ", posted: 0, fixed: 0},{content: "い", posted: 0, fixed: 1}])}>マイページ</h1>
                <p>ログイン中のユーザー：{user?.email}{send}</p>
                <button onClick={logout}>ログアウト</button>
              </header>
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
                </> : null}
              </footer>
            </>
          )}
        </>
      )}
    </>
  );
};

export default Mypage;