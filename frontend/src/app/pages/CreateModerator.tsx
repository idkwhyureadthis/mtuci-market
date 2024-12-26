import React, { useEffect, useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import { useNavigate } from 'react-router-dom';

export function CreateModerator() {
  const navigate = useNavigate() 
  const [formData, setFormData] = useState({
    login: '',
    password: '',
    name: '',
    telegram: '',
  });

  type UserData = {
    name: string,
    new_tokens: {
      access: string
      refresh: string
    },
    dorm_number: string,
    role: string,
}

  const [loginOccupied, setLoginOccupied] = useState(false)

  useEffect(() => {
    const tokens = {
      access: sessionStorage.getItem("access"),
      refresh: localStorage.getItem("refresh")
    }
    if (tokens.access === null && tokens.refresh === null){
      navigate("/sign-in")
    }
    checkTokens(tokens);
}, []);


const checkTokens = async (tokens: any) => {
  try{
    const response : AxiosResponse = await axios.post("http://localhost:8082/verify", tokens)
    const userData : UserData = response.data
    if (userData.new_tokens !== null && userData.new_tokens !== undefined){
      localStorage.setItem("refresh", userData.new_tokens.refresh)
      sessionStorage.setItem("access", userData.new_tokens.access)
    }
    if (userData.role != "admin"){
      navigate("/")
    }
  } catch (e)
  {
    console.error(e)
    localStorage.clear()
    sessionStorage.clear()
    navigate("/")
  }
}


  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try{
        const response: AxiosResponse = await axios.post("http://localhost:8082/create_moderator", formData)
        if (response.status != 200) {
          console.log("ватафак")
        }
    } catch (e) {
        console.log(e)
        setLoginOccupied(true)
    }
    navigate("/profile")
  };

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center py-5 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8 bg-white p-8 rounded-lg shadow">
        <div>
          <h2 className="text-center text-2xl font-bold text-gray-900">
            Регистрация Модератора
          </h2>
        </div>
        <form className="mt-8 space-y-3" onSubmit={handleSubmit}>
          <div className="space-y-4">
            <div>
              <label htmlFor="username" className="block text-sm font-medium text-gray-700">
                Логин
              </label>
              <input
                id="username"
                type="text"
                required
                className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-1"
                value={formData.login}
                onChange={(e) => setFormData(prev => ({ ...prev, login: e.target.value }))}
              />
            </div>
            <div>
              <label htmlFor="password" className="block text-sm font-medium text-gray-700">
                Пароль
              </label>
              <input
                id="password"
                type="password"
                required
                className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-1"
                value={formData.password}
                onChange={(e) => setFormData(prev => ({ ...prev, password: e.target.value }))}
              />
            </div>
            <div>
              <label htmlFor="name" className="block text-sm font-medium text-gray-700">
                Имя
              </label>
              <input
                id="name"
                type="text"
                required
                className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-1"
                value={formData.name}
                onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
              />
            </div>
            <div>
              <label htmlFor="telegram" className="block text-sm font-medium text-gray-700">
                Telegram
              </label>
              <div className="mt-1 flex rounded-md shadow-sm">
                <span className="inline-flex items-center px-3 rounded-l-md border border-r-0 border-gray-300 bg-gray-50 text-gray-500 text-sm">
                  @
                </span>
                <input
                  id="telegram"
                  type="text"
                  required
                  placeholder="username"
                  className="flex-1 block w-full rounded-none rounded-r-md border border-gray-300 px-3 py-1"
                  value={formData.telegram}
                  onChange={(e) => setFormData(prev => ({ ...prev, telegram: e.target.value }))}
                />
              </div>
            </div>
          </div>
          <button
            type="submit"
            className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-800 disabled:hover:none disabled:bg-blue-300"
            disabled = {(formData.login == "" || formData.telegram == "" || formData.name == "" || formData.password == "")}
          >
            Зарегистрироваться
          </button>

          
          {loginOccupied ? 
          <div className="w-full bg-white shadow-lg rounded-lg relative">
                <div className="bg-red-400 text-primary-50 rounded-md px-5 py-2 text-center">
                  <p className="font-title">Ошибка: Логин уже занят</p>
                </div>
          </div> : ''
          }
        </form>
      </div>
    </div>
  );
}