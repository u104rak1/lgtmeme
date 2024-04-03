export const toSnakeCase = (str: string): string => {
  return str.replace(
    /([A-Z])/g,
    (char, index) => `${index > 0 ? "_" : ""}${char.toLowerCase()}`
  );
};
