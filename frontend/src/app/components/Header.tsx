import { Plus, LogOut, UserPlus, Sparkles } from 'lucide-react';
import { useNavigate} from "react-router-dom"

interface HeaderProps{
  name: string
  status: string
}

export function Header({name, status}: HeaderProps) {
  const navigate = useNavigate()
  return (
    <header className="bg-white shadow-sm z-10">
      <div className="max-w-7xl mx-auto px-4 py-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center">
          <a href='/'><h1 className="text-2xl font-bold text-gray-900">Обосраная жопа</h1></a>
          <div className="flex gap-10">
            {status == "admin" ?
            <button className="bg-blue-600 px-6 py-2 rounded-xl flex items-center gap-6 text-blue-50 hover:bg-blue-700" onClick={() => {
              navigate("/create_moderator")
            }}>
              <UserPlus size={20}/>Создать модератора
            </button> : ""}
            {status == "moderator" ?
            <button className="bg-blue-600 px-6 py-2 rounded-xl flex items-center gap-6 text-blue-50 hover:bg-blue-700" onClick={() => {
              navigate("/moderator")
            }}>
              <Sparkles size={20}/>Модерация карточек
            </button> : ""}
            <button className="bg-blue-600 text-blue-50 px-4 py-2 rounded-xl flex items-center gap-5 hover:bg-blue-700"
            onClick = {
              () => {navigate("/create")}
            }>
              <Plus size={20} />
              Продать вещь
            </button>
            <a href="/profile"><div className="flex align-middle my-2">
              {name}
              </div>
            </a>
              <div>
                <button className="flex align-middle rounded-xl items-center py-2 px-4 bg-gray-500 hover:bg-gray-400 text-neutral-100 hover:text-neutral-50"
                onClick={ () => {
                    localStorage.clear()
                    sessionStorage.clear()
                    navigate('/sign-in')
                  }
                }
                >
                    Выйти из аккаунта <LogOut size={20} className='mx-2'></LogOut>
                </button>
              </div>
          </div>
        </div>
      </div>
    </header>
  );
}