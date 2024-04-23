import { css, cva } from "@@/styled-system/css";

type Props = {
  css?: string;
  tabs: {
    id: string;
    label: string;
  }[];
  value: string;
  onClick: (value: string) => void;
};

const Tabs = ({ css, tabs, value, onClick }: Props) => {
  const handleClick = (value: string) => onClick(value);
  return (
    <div className={css}>
      <div className={tabsCss}>
        {tabs.map((tab) => (
          <div
            key={tab.id}
            className={tabRecipe({ active: tab.id === value })}
            onClick={() => handleClick(tab.id)}
          >
            {tab.label}
          </div>
        ))}
      </div>
    </div>
  );
};

const tabsCss = css({ display: "flex", width: "100%", boxShadow: "lg" });
const tabRecipe = cva({
  base: {
    flex: 1,
    padding: "3",
    cursor: "pointer",
    bgColor: "WHITE",
    color: "BLACK",
    textAlign: "center",
    opacity: "0.8",
  },
  variants: {
    active: {
      true: {
        fontWeight: "bold",
        borderBottom: "2px solid",
        borderBottomColor: "BLACK",
        opacity: "1",
      },
    },
  },
});

export default Tabs;
