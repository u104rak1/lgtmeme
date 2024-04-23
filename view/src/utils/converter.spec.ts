import { toSnakeCase } from "@/utils/converter";

describe("toSnakeCase", () => {
  it("converts simple camelCase to snake_case", () => {
    expect(toSnakeCase("camelCase")).toBe("camel_case");
  });
  it("handles strings without any uppercase letters", () => {
    expect(toSnakeCase("lowercase")).toBe("lowercase");
  });
  it("handles strings that are already in snake_case", () => {
    expect(toSnakeCase("already_snake_case")).toBe("already_snake_case");
  });
  it("handles mixed strings with numbers", () => {
    expect(toSnakeCase("version1Point2")).toBe("version1_point2");
  });
  it("handles strings with leading and trailing underscores", () => {
    expect(toSnakeCase("_leadingCamelCase_")).toBe("_leading_camel_case_");
  });
});
