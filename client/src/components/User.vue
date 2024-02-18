<script setup>
import { ref, unref } from 'vue';
import axios from 'axios';

const state = ref('login');

const user = ref();

const input = ref({
  email: '',
  password: '',
});

const error = ref();

const signup = async () => {
  error.value = null;
  try {
    await axios.post('http://localhost:3001/auth/register', unref(input));
    login();
  } catch (e) {
    error.value = `${e.response.status}: ${e.response.data}`;
  }
};

const login = async () => {
  error.value = null;
  try {
    const res = await axios.post(
      'http://localhost:3001/auth/login',
      unref(input)
    );
    user.value = res.data;

    state.value = 'notif';
    connect();
  } catch (e) {
    error.value = `${e.response.status}: ${e.response.data}`;
  }
};

const messages = ref([]);

const connect = async () => {
  document.cookie = 'X-Authorization=' + user.value.token + '; path=/';

  var socket = new WebSocket('ws://localhost:3001/ws/user/notifications');
  socket.addEventListener('open', (event) => {
    console.log(user.value.id + ' connected');
  });

  socket.addEventListener('message', (event) => {
    messages.value.push({
      message: event.data,
      time: new Date().toLocaleTimeString(),
    });
  });
};

const sendMessage = async () => {
  axios.post(`http://localhost:3001/user/${user.value.id}/notifications`, {
    message:
      'Some random message at: ' + Math.floor(new Date().getTime() / 1000),
  });
};
</script>

<template>
  <div class="user-box">
    <div class="inputs" v-if="state === 'login'">
      <label>Email:</label>
      <input type="email" v-model="input.email" />
      <label>Password:</label>
      <input type="password" v-model="input.password" />

      <div class="err" v-if="error">{{ error }}</div>

      <button @click="login">LOGIN</button>
      <button @click="signup">SIGNUP</button>
    </div>

    <div class="messages" v-if="state === 'notif'">
      <p>Messages:</p>
      <div class="msg" v-for="m of messages">{{ m.time }}: {{ m.message }}</div>

      <button @click="sendMessage">Send random message to user</button>
    </div>
  </div>
</template>

<style>
.user-box {
  border: 1px solid black;
  padding: 5%;
  width: 400px;
}

.inputs,
.messages {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.msg {
  width: 100%;
}

.err {
  color: red;
}
</style>
