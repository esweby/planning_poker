import cl from "./Select.module.css";

interface SelectProps {
  options: { value?: string; display: string }[];
  className?: string;
  value: string;
  onChange: (value: any) => void;
}

const Select = ({ options, className, value, onChange }: SelectProps) => {
  let classes = cl.select;
  if (className) classes += ` ${className}`;

  return (
    <select
      className={classes}
      value={value}
      onChange={(e) => onChange(e.target.value as any)}
    >
      {options.map((opt) => (
        <option value={opt?.value || ""}>{opt.display}</option>
      ))}
    </select>
  );
};

export default Select;
