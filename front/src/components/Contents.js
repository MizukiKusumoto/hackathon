import React, { useState, useEffect } from "react";

const Contents = () => {
  const [messages, setMessages] = useState([]);
  const [nowChannel, setNow] = useState({});
  const [postedMessage, setPost] = useState("");
  const [channels, setChannels] = useState([]);
  const [updating, setUpdating] = useState({});
  const [updatedMessage, setUpdated] = useState("");

  const getChannels = async () => {
    try {
        const response = await fetch("https://myhackathon-7iuhbg7yzq-uc.a.run.app/channel");
        if (!response.ok) {
            throw Error(`Failed to fetch users: ${response.status}`);
        }
        const data = await response.json();
        setChannels(data);
    } catch(error) {
        console.error( "エラー：", error );
    }
  };

  useEffect(() => {
    getChannels();
  }, []);

  const getMessages = async () => {
    if (nowChannel === null){
        setMessages([]);
        return;
    }
    try {
        const url = "https://myhackathon-7iuhbg7yzq-uc.a.run.app/message?channel_id=" + nowChannel.id;
        const response = await fetch(url);
        if (!response.ok) {
            throw Error(`Failed to fetch users: ${response.status}`);
        }
        const data = await response.json();
        setMessages(data);
    } catch(error) {
        console.error( "エラー：", error );
    }
  };

  useEffect(() => {
    getMessages();
  }, [nowChannel]);

  const postMessage = async (message, channelId) => {
    try {
      const response = await fetch("https://myhackathon-7iuhbg7yzq-uc.a.run.app/message", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ content: message, channel_id: channelId }),
      });
      if (!response.ok) {
        throw new Error(`Failed to add message: ${response.status}`);
      }
    } catch (error) {
      console.error("エラー：", message);
    }
    setPost("");
    getMessages();
  };

  const updateMessage = async (id, message) => {
    try {
      const response = await fetch("https://myhackathon-7iuhbg7yzq-uc.a.run.app/message", {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ id: id, content: message }),
      });
      if (!response.ok) {
        throw new Error(`Failed to update message: ${response.status}`);
      }
    } catch (error) {
      let message
      if (error instanceof Error) message = error.message
      else message = String(error)
      console.error("エラー：", message);
    }
    setUpdated("");
    setUpdating({});
    getMessages();
  };

  const deleteMessage = async (id) => {
    try {
      const response = await fetch(`https://myhackathon-7iuhbg7yzq-uc.a.run.app/message?id=${id}`, {
        method: "DELETE",
      });
      if (!response.ok) {
        throw new Error(`Failed to delete message: ${response.status}`);
      }
    } catch (error) {
      let message
      if (error instanceof Error) message = error.message
      else message = String(error)
      console.error("エラー：", message);
    }
    getMessages();
  };

  return (
    <>
        <main>
            <ul>
                {!(nowChannel.id == null) ? <>
                    {!(messages.length == 0) ? <>
                        {messages.map((message) => (<>
                            <li>
                                {message.content}{!(message.created_at == message.modified_at) ? <>(編集済み)</> : null}
                                <button onClick={() => {setUpdating({id: message.id, now: true }); setUpdated(message.content);}}>編集</button>
                                <button onClick={() => deleteMessage(message.id)}>削除</button>
                            </li>
                            {updating.now ? <>{updating.id == message.id ? <><textarea
                                onChange={(e) => setUpdated(e.target.value)}
                                >{updatedMessage}</textarea>
                                <button onClick={() =>updateMessage(updating.id, updatedMessage)}>修正</button>
                                </> : null
                            }</> : null}
                        </>))} </> : <li>会話を始めましょう！</li>
                    }</> : null
                }
            </ul>
        </main>
        <aside>
            {!(nowChannel.id == null) ? <p>あなたの現在のチャンネル：{nowChannel.name}</p> : <p>チャンネルを選択してください。</p>}
            {channels.map((channel) => (
                <button onClick={() => setNow({ id: channel.id, name: channel.name })}>{channel.name}</button>
            ))}
        </aside>
        <footer>
            {!(nowChannel.id == null) ? <><textarea
                onChange={(e) => setPost(e.target.value)}
                >{postedMessage}</textarea>
                <button onClick={() =>postMessage(postedMessage, nowChannel.id)}>送信</button>
                </> : null
            }
        </footer>
    </>
  );
};

export default Contents;