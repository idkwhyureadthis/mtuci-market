import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import axios, { AxiosResponse } from "axios";
import { Header } from "../components/Header";
import { Product } from "../types";
import { ProductCard } from "../components/ProductCard";


const UserProfile = () => {
     
    type UserData = {
        id: number
        name: string,
        new_tokens: {
          access: string
          refresh: string
        },
        dorm_number: string,
        role: string,
    }

    type CardsData = {
        cards: Product[]
    }

    const [products, setProducts] = useState([] as Product[])
    const [userName, setUserName] = useState("")
    const [dormNumber, setDormNumber] = useState("")
    const [role, setRole] = useState("user")

    const navigate = useNavigate()


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

    const getProducts = async (id: number) => {
        try{
            const response : AxiosResponse = await axios.get("http://localhost:8082/products_of?id="+id)
            const data: CardsData = response.data
            setProducts(data.cards)
        } catch (e){
            console.log(e)
        }
    }

    const sleep = (ms:number) => {
        return new Promise(resolve => setTimeout(resolve, ms));
    }

    const checkTokens = async (tokens: any) => {
        try{
          const response : AxiosResponse = await axios.post("http://localhost:8082/verify", tokens)
          const userData : UserData = response.data
          setUserName(userData.name)
          setDormNumber(userData.dorm_number)
          setRole(userData.role)
          console.log(userData)
          sleep(100).then(() => getProducts(userData.id))
          if (userData.new_tokens !== null && userData.new_tokens !== undefined){
            localStorage.setItem("refresh", userData.new_tokens.refresh)
            sessionStorage.setItem("access", userData.new_tokens.access)
          }
        } catch (e)
        {
          console.error(e)
          navigate("/sign-in")
        }
      }
    
    const deleteCard = async (id: string) => {
        setProducts(products.filter((product) => product.id != id))
        try{
            axios.post("http://localhost:8082/delete_card/" + id)
        } catch(e){
            console.log(e)
        }
    }

    return(
        <div className="min-h-screen bg-gray-50">
        <Header name={userName} status={role}/>
        <main className="max-w-7xl mx-auto px-4 py-8 sm:px-6 lg:px-8">
            <div> 
                <h1 className="text-center text-3xl font-title text-neutral-950 mb-6">Профиль пользователя:</h1>
              <section className="mb-10">
                <div className="text-center">
                  <h2 className="text-xl font-semibold text-neutral-950">{userName}</h2>
                  {role == "user" ? <p className="text-neutral-600">Студент общежития №{dormNumber}</p> : ""}
                  {role == "moderator" ? <p className="text-neutral-600">Модератор{dormNumber}</p> : ""}
                  {role == "admin" ? <p className="text-neutral-600">Администратор{dormNumber}</p> : ""}
                  <h1 className="text-3xl mt-2"> Вещи на продаже:</h1>
                </div>
              </section>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {products.map((product) => (
                  <div className="flex flex-col justify-center text-center">
                    <ProductCard product={product} place='profile'/>
                    <div className={product.status == "accepted" ? "z-30 -mt-20 mb-2 text-green-400" : product.status == "rejected" ? "z-30 -mt-20 mb-2 text-red-400" : "z-30 -mt-20 mb-2 text-neutral-700"}>{product.status == "accepted" ? "Принято" : ""}{product.status == "rejected" ? "Отклонено" : ""}{product.status == "on moderation" ? "На модерации" : ""}</div>
                    <button className="mx-16 py-2 bg-red-400 hover:bg-red-500 text-red-50 rounded-full" onClick={() => {deleteCard(product.id)}}>Удалить</button>
                  </div>
                ))
            }
            </div>
        </main>
        </div>
    )
}

export default UserProfile;