import React, { useEffect, useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import { useNavigate } from 'react-router-dom';

export function SignUpPage() {
  const navigate = useNavigate() 
  const [formData, setFormData] = useState({
    login: '',
    password: '',
    name: '',
    room: '',
    dorm_number: '1',
    telegram: '',
  });

  useEffect(() =>{
  if (localStorage.getItem("refresh") != null && sessionStorage.getItem("access") != null){
    navigate("/")
  }
  }, [])

  const [loginOccupied, setLoginOccupied] = useState(false)

  type tokensResponse = {
      access: string,
      refresh: string,
  }

  const handleSubmit = async (e: React.FormEvent) => {
    var responseData: tokensResponse | undefined = undefined
    e.preventDefault();
    try{
        const response: AxiosResponse = await axios.post("http://localhost:8082/signup", formData)
        responseData = response.data
    } catch {
        setLoginOccupied(true)
    }finally {
      if (responseData != undefined) {
        localStorage.setItem("refresh", responseData.refresh)
        sessionStorage.setItem("access", responseData.access)
        navigate("/profile")
      }
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center py-5 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8 bg-white p-8 rounded-lg shadow">
        <div>
          <h2 className="text-center text-2xl font-bold text-gray-900">
            Регистрация
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
            <div>
              <label htmlFor="dormitory" className="block text-sm font-medium text-gray-700">
                Общежитие
              </label>
              <select
                id="dormitory"
                required
                className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-1"
                value={formData.dorm_number}
                onChange={(e) => setFormData(prev => ({ ...prev, dorm_number: e.target.value as '1' | '3' | '4' }))}
              >
                <option value="1">№1</option>
                <option value="3">№3</option>
                <option value="4">№4</option>
              </select>
            </div>
            <div>
              <label htmlFor="room" className="block text-sm font-medium text-gray-700">
                Комната
              </label>
              <input
                id="room"
                type="text"
                required
                className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-1"
                value={formData.room}
                onChange={(e) => setFormData(prev => ({ ...prev, room: e.target.value }))}
              />
            </div>
          </div>
          <button
            type="submit"
            className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-800 disabled:hover:none disabled:bg-blue-300"
            disabled = {(formData.login == "" || formData.room == "" || formData.telegram == "" || formData.name == "" || formData.password == "")}
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