import { Link, useNavigate } from "react-router-dom"
import Textbox from "./components/Textbox"
import Button from "./components/Button"
import { useForm } from "react-hook-form"
import { UserRequest, createUser } from "./api"
import { useMutation } from "@tanstack/react-query"
import { useAuthAction } from "../../store"
import toast from "react-hot-toast"

function Login() {
  const {
    formState: { isDirty, isSubmitting },
    handleSubmit,
    register,
  } = useForm<UserRequest>()

  const navigate = useNavigate()

  const setId = useAuthAction().setId
  const setUsername = useAuthAction().setUsername

  const mutation = useMutation({
    mutationFn: createUser
  })

  const onSubmit = async (req: UserRequest) => {
    try {
      const res = await mutation.mutateAsync(req)
      setId(res.data._id)
      setUsername(res.data.UserName)
      toast.success("Berhasil membuat akun")
      navigate("/")
    } catch (err) {
      toast.error("Gagal membuat akun")
    }
  };
  
  return (
    <div className="w-full h-[100vh] overflow-hidden bg-zinc-900 flex flex-col items-center justify-center font-sono">
      <span className="block my-4 text-3xl font-bold text-white">Register</span>
      <form className="w-[400px] max-w-[80%]" onSubmit={handleSubmit(onSubmit)}>
        <Textbox label="username" className="my-4 p-4 focus:outline-white" name="username" register={register} required={true}/>
        <Textbox label="password" type="password" className="my-4 focus:outline-white" name="password" register={register} required={true}/>
        <div className="flex justify-center">
          <Button label="submit" type="submit" className="bg-yellow-200 mx-4 cursor-pointer" disabled={!isDirty} loading={isSubmitting}/>
          <Link to="/auth/login">
            <Button label="login" className="bg-indigo-600 text-white"/>
          </Link>
        </div>
      </form>
    </div>
  )
}

export default Login