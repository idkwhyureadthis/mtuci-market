import { useState, useEffect } from 'react';
import { Product } from '../types';
import { ImageGallery } from '../components/ImageGallery';
import { useNavigate } from 'react-router-dom';
import axios, { AxiosResponse } from 'axios';
import { Header } from '../components/Header';


export function ModeratorPage() {

  const navigate = useNavigate()

  const handleApprove = async (productId: string) => {
    try{
      axios.post("http://localhost:8082/accept/" + productId)
    } catch (e){
      console.log(e)
    } finally{
      window.location.reload()
    }
  };

  const handleReject = (productId: string) => {
    try{
      axios.post("http://localhost:8082/reject/" + productId)
    } catch (e){
      console.log(e)
    } finally{
      window.location.reload()
    }
  };

  const [product, setProduct] = useState({} as Product);
  const [userName, setUserName] = useState("")
  const [status, setStatus] = useState("user")



  type UserData = {
    name: string,
    new_tokens: {
      access: string
      refresh: string
    },
    dorm_number: string,
    role: string,
}
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

const getCard = async () => {
  const resp : AxiosResponse = await axios.get("http://localhost:8082/on_moderation")
  const data : Product = resp.data
  if (data.photos === null || data.photos === undefined){
    data.photos = []
  }
  console.log(data)
  setProduct(data)
}

const checkTokens = async (tokens: any) => {
  try{
    const response : AxiosResponse = await axios.post("http://localhost:8082/verify", tokens)
    const userData : UserData = response.data
    setUserName(userData.name)
    setStatus(userData.role)
    if (userData.new_tokens !== null && userData.new_tokens !== undefined){
      localStorage.setItem("refresh", userData.new_tokens.refresh)
      sessionStorage.setItem("access", userData.new_tokens.access)
    }
    if (userData.role != "moderator"){
      navigate("/")
    }
    getCard()
  } catch (e)
  {
    console.error(e)
    localStorage.clear()
    sessionStorage.clear()
    navigate("/")
  }
}

  return (
    <div>
    <Header name={userName} status={status}></Header>
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <h1 className="text-3xl mt-10 font-bold text-gray-900 mb-10 text-center">
          Модерация объявлений
        </h1>
        <div className="grid grid-cols-1 gap-6">
          {product.id != "0" ? 
            <div key={product.id} className="bg-white rounded-lg shadow p-6 h">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <ImageGallery photos={product.photos ?? []} title={product.name} />
                </div>
                <div>
                  <h2 className="text-xl font-bold mb-2">{product.name}</h2>
                  <p className="text-gray-600 mb-4">{product.description}</p>
                  <div className="mb-4">
                    <span className="text-xl font-bold text-blue-600">
                      {product.price} ₽
                    </span>
                  </div>
                  <div className="mb-4">
                    <p className="text-sm text-gray-600">
                      Продавец: {product.creator_name}
                    </p>
                    <p className="text-sm text-gray-600">
                      Общежитие №{product.dorm_number}, комната {product.room}
                    </p>
                  </div>
                  <div className="flex gap-4">
                    <button
                      onClick={() => {
                        handleApprove(product.id)}}
                      className="flex-1 py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-600 hover:bg-green-700"
                    >
                      Одобрить
                    </button>
                    <button
                      onClick={() => {
                        handleReject(product.id)}}
                      className="flex-1 py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-red-500 hover:bg-red-900"
                    >
                      Отклонить
                    </button>
                  </div>
                </div>
              </div>
            </div>
            : ""}
          <div className='text-center text-5xl mt-48'>{product.id == "0" ? "Вы рассмотрели все карточки, нуждающиеся в модерации": ""}</div>
        </div>
      </div>
    </div>
    </div>
  );
}