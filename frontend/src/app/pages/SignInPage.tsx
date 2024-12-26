import React, { useState } from 'react';
import axios, { AxiosResponse } from 'axios';
import { useNavigate } from 'react-router-dom';
import { useEffect } from 'react';

export function SignInPage() {
  const navigate = useNavigate() 
  const [formData, setFormData] = useState({
    login: '',
    password: '',
  });

  type tokensResponse = {
    access: string,
    refresh: string,
  }
  
  useEffect(() =>{
    if (localStorage.getItem("refresh") != null && sessionStorage.getItem("access") != null){
      navigate("/")
    }
  },[])


  const [wrongData, setWrongData] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    var responseData: tokensResponse | undefined = undefined
    e.preventDefault();
    try{
        const response: AxiosResponse = await axios.post("http://localhost:8082/login", formData)
        responseData = response.data
        if (responseData != undefined) {
          localStorage.setItem("refresh", responseData.refresh)
          sessionStorage.setItem("access", responseData.access)
          navigate("/")
        }
    } catch (e) {
      setWrongData(true)
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8 bg-white p-8 rounded-lg shadow">
        <div>
          <h2 className="text-center text-3xl font-bold text-gray-900">
            Вход
          </h2>
        </div>
        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="space-y-4">
            <div>
              <label htmlFor="username" className="block text-sm font-medium text-gray-700">
                Логин
              </label>
              <input
                id="username"
                type="text"
                required
                className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2"
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
                className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2"
                value={formData.password}
                onChange={(e) => setFormData(prev => ({ ...prev, password: e.target.value }))}
              />
            </div>
          </div>
          <button
            type="submit"
            className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-800 disabled:hover:none disabled:bg-blue-300"
            disabled = {(formData.password == "" || formData.login == "" )}
          >
            Войти
          </button>
          {wrongData == true ? 
          <div className="w-full bg-white shadow-lg rounded-lg relative">
                <div className="bg-red-400 text-primary-50 rounded-md px-5 py-2 text-center">
                  <p className="font-title">Ошибка: Введённые данные неверны</p>
                </div>
          </div> : ''
          }
        </form>
      </div>
    </div>
  );
}