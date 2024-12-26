import React, { useState, useEffect } from 'react';
import axios, {AxiosResponse} from "axios";
import {useNavigate} from "react-router-dom";
import { Header } from '../components/Header';

  
type UserData = {
  name: string,
  new_tokens: {
    access: string
    refresh: string
  },
  dorm_number: string,
  role: string,
}


type RequestTokens = {
  tokens: {
    access: string,
    refresh: string,
  }
}


export function CreateProductPage() {


  useEffect(() => {
    const tokens = {
      access: sessionStorage.getItem("access"),
      refresh: localStorage.getItem("refresh")
    }
    if (tokens.access === null || tokens.refresh === null){
      navigate("/sign-in")
    }
    checkTokens(tokens);
  }, []);
  


  const [formData, setFormData] = useState({
    name: '',
    description: '',
    price: '',
  });

  const [userName, setUserName] = useState("")

  const navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const form = new FormData()
    images.forEach(element => {
        form.append("images", element)
    });
    form.append("refresh", localStorage.getItem("refresh") ?? "")
    form.append("access", sessionStorage.getItem("access") ?? "")
    form.append("name", formData.name)
    form.append("description", formData.description)
    form.append("price", formData.price)
    try{
      const response : AxiosResponse = await axios.post("http://localhost:8082/create_card", form)
      const tokens : RequestTokens = response.data
      if (tokens.tokens !== null && tokens.tokens !== undefined && tokens.tokens.access != ""){
        localStorage.setItem("refresh", tokens.tokens.refresh)
        sessionStorage.setItem("access", tokens.tokens.access)
      }
      setFormData({name: '', description: '', price: ''})
      setPreviews([])
      setImages([])
    } catch (e){
      console.log(e)
    }
  };

  const checkTokens = async (tokens: any) => {
    try{
      const response : AxiosResponse = await axios.post("http://localhost:8082/verify", tokens)
      const userData : UserData = response.data
      setUserName(userData.name)
      setStatus(userData.role)
      if (userData.new_tokens !== null && userData.new_tokens !== undefined && userData.new_tokens.access != ""){
        localStorage.setItem("refresh", userData.new_tokens.refresh)
        sessionStorage.setItem("access", userData.new_tokens.access)
      }
    } catch (e){
      console.error(e)
      navigate("/sign-in")
    }
  }


  const [images, setImages] = useState<File[]>([]);
  const [previews, setPreviews] = useState<string[]>([]);
  const [status, setStatus] = useState("user")

  const handleImageChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files) {
      const files = Array.from(event.target.files);
      setImages(files);
      const previews = files.map((file) => URL.createObjectURL(file));
      setPreviews(previews);
    }
  };

  const handleRemoveImage = (index: number) => {
    const newImages = images.filter((_, i) => i !== index);
    const newPreviews = previews.filter((_, i) => i !== index);
    setImages(newImages);
    setPreviews(newPreviews);
  };


  const tokens = {
    access: sessionStorage.getItem("access"),
    refresh: localStorage.getItem("refresh")
  }
  
  if (tokens.access === null || tokens.refresh === null){
    navigate("/sign-in")
  }



  return (
    <div>
    <Header name={userName} status={status}/>
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-2xl font-bold text-gray-900 mb-6">
            Создать объявление
          </h2>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label htmlFor="title" className="block text-sm font-medium text-gray-700">
                Название
              </label>
              <input
                id="title"
                type="text"
                required
                className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2"
                value={formData.name}
                onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
              />
            </div>
            <div>
              <label htmlFor="description" className="block text-sm font-medium text-gray-700">
                Описание
              </label>
              <textarea
                id="description"
                rows={4}
                required
                className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2"
                value={formData.description}
                onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
              />
            </div>
            <div>
              <label htmlFor="price" className="block text-sm font-medium text-gray-700">
                Цена (₽)
              </label>
              <input
                id="price"
                type="number"
                min="0"
                required
                className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2"
                value={formData.price}
                onChange={(e) => setFormData(prev => ({ ...prev, price: e.target.value }))}
              />
            </div>
            <div className="flex flex-col gap-4 p-4 border border-gray-200 rounded-lg">
              <input
                type="file"
                multiple
                accept=".jpg, .jpeg, .png"
                onChange={handleImageChange}
                className="block w-full text-sm text-gray-500 mt-2 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
              />
              <div className="flex flex-wrap gap-4">
                {previews.map((preview, index) => (
                  <div key={index} className="relative">
                    <img
                      src={preview}
                      alt="Uploaded Image"
                      className="w-16 h-16 object-cover rounded-lg"
                    />
                    <button
                      onClick={() => handleRemoveImage(index)}
                      className="absolute top-0 right-0 py-0 px-2 bg-red-500 rounded-full text-white hover:bg-red-700 align-middle text-center"
                    >
                      ×
                    </button>
                  </div>
                ))}
              </div>
            </div>
            <button
              type="submit"
              className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-800 disabled:hover:none disabled:bg-blue-300"
              disabled = {formData.price == "0" || formData.price == '' || formData.description == "" || formData.name == ""}
            >
              Создать объявление
            </button>
          </form>
        </div>
      </div>
    </div>
    </div>
  );
}