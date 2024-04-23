export const toSnakeCase = (str: string): string => {
  return str.replace(/([A-Z])/g, (_match, char) => `_${char.toLowerCase()}`);
};
