import React, { useState, useEffect } from "react";
import { onAuthStateChanged, signOut } from "firebase/auth";
import { auth } from "../firebase.js";
import { useNavigate, Navigate } from "react-router-dom";
import Contents from "./Contents.js";

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
    <div className="mypageContainer">
      {!loading && (
        <>
          {!user ? (
            <Navigate to={`/login/`} />
          ) : (
            <>
              <header>
                <div>
                  <h1 onClick={() => setMessages([{content: "あ", posted: 0, fixed: 0},{content: "い", posted: 0, fixed: 1}])}>マイページ</h1>
                  <p>ログイン中のユーザー：{user?.email}{send}</p>
                </div>
                <button onClick={logout}>ログアウト</button>
              </header>
              <Contents />
            </>
          )}
        </>
      )}
    </div>
  );
};

export default Mypage;