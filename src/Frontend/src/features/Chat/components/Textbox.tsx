import { HistoryRequest } from "../api"
import clsx from "clsx"
import { InputHTMLAttributes } from "react"
import { UseFormRegister} from "react-hook-form"

type TextboxProps = InputHTMLAttributes<HTMLInputElement> & {
    label?: string
    className?: string
    name: "userId" | "name",
    register: UseFormRegister<HistoryRequest> 
}

function Textbox({
    label,
    name,
    required,
    disabled,
    className,
    register,
    ...props
}: TextboxProps) {
  return (
    <input
        {...register(name, { required: required} )}
        disabled={disabled}
        placeholder={label}
        {...props}
        className={clsx("w-full outline-none p-4 bg-zinc-950 flex justify-between items-center rounded-3xl text-white", className)}
    />
  )
}

export default Textbox