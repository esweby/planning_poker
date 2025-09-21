import cl from "./Input.module.css";

interface InputProps {
  name: string;
  placeholder: string;
  value: string;
  onChange: (name: string) => void;
  className?: string;
}

const Input = ({
  name,
  placeholder,
  value,
  onChange,
  className,
}: InputProps) => {
  let classes = cl.input;
  if (className) classes += ` ${className}`;
  return (
    <input
      type="text"
      name={name}
      placeholder={placeholder}
      className={classes}
      value={value}
      onChange={(e) => onChange(e.target.value)}
    />
  );
};

export default Input;
