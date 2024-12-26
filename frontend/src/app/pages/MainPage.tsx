import { Search } from 'lucide-react';
import { ProductCard } from '../components/ProductCard';
import { DormitoryFilter } from '../components/DormitoryFilter';
import { Pagination } from '../components/Pagination';
import { useProducts } from '../hooks/useProducts';
import { Header } from '../components/Header';
import axios, {AxiosResponse} from "axios";
import {useNavigate} from "react-router-dom"
import {useState, useEffect} from 'react';
import {Product} from '../types';
 
type UserData = {
    name: string,
    new_tokens: {
      access: string
      refresh: string
    },
    dorm_number: string,
    role: string,
}

type GetCardsData = {
  cards: Product[]
}


export function MainPage(){

  const [parsedCards, setParsedCards] = useState([] as Product[])
  var cards: GetCardsData;

  useEffect(() => {
      const tokens = {
        access: sessionStorage.getItem("access"),
        refresh: localStorage.getItem("refresh")
      }
      if (tokens.access === null && tokens.refresh === null){
        navigate("/sign-in")
      }
      checkTokens(tokens);

      getCards()
  }, []);

  const getCards = async () => {
    const response : AxiosResponse = await axios.get("http://localhost:8082/get_cards")
    cards = response.data
    console.log(cards)

    setParsedCards(cards.cards)
  }



  const navigate = useNavigate()
  const [userName, setUserName] = useState("")
  const [status, setStatus] = useState("user")
  const {
    products,
    currentPage,
    totalPages,
    setCurrentPage,
    selectedDormitory,
    setSelectedDormitory,
    searchQuery,
    setSearchQuery,
  } = useProducts(parsedCards);


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
    } catch (e)
    {
      console.error(e)
      localStorage.clear()
      sessionStorage.clear()
      navigate("/sign-in")
    }
  }


  return (
    <div className="min-h-screen bg-gray-50">
      <Header name={userName} status={status}></Header>
      <main className="max-w-7xl mx-auto px-4 py-8 sm:px-6 lg:px-8">
        <div className="mb-8">
          <div className="relative">
            <input
              type="text"
              placeholder="Поиск вещей..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full px-4 py-2 pl-10 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
            <Search className="absolute left-3 top-2.5 text-gray-400" size={20} />
          </div>
        </div>

        <DormitoryFilter
          selectedDormitory={selectedDormitory}
          onSelect={setSelectedDormitory}
        />

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {products.map((product) => (
            <ProductCard key={product.id} product={product} place='mainpage'/>
          ))}
        </div>

        <Pagination
          currentPage={currentPage}
          totalPages={totalPages}
          onPageChange={setCurrentPage}
        />
      </main>
    </div>
  );
}