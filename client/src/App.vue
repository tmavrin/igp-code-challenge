<script setup>
import { ref } from 'vue';
import User from './components/User.vue';

const users = ref([]);

const addUser = () => {
  users.value.push({ id: Math.floor(new Date().getTime() / 1000) });
};
</script>

<template>
  <div class="container">
    <User v-for="user of users" :key="user.id" />
    <button class="login" @click="addUser">Add User</button>
  </div>
</template>

<style scoped>
.container {
  margin: auto;
  width: 100%;
  height: 100%;
  max-width: 1200px;
  align-items: center;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

button.login {
  font-size: 1.1em;
}
</style>
