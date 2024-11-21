import pandas as pd

def find_min_max(file_path):
    df = pd.read_csv(file_path, delimiter='|', na_values=[''], keep_default_na=True)

    numeric_columns = df.select_dtypes(include=['float64', 'int64']).columns

    min_values = {}
    max_values = {}

    for column in numeric_columns:
        lower_quantile = df[column].quantile(0.10)
        upper_quantile = df[column].quantile(0.90)
        filtered_df = df[(df[column] >= lower_quantile) & (df[column] <= upper_quantile)]

        min_id = filtered_df[column].idxmin()
        max_id = filtered_df[column].idxmax()

        min_values[column] = {
            'id': df.loc[min_id, 'id'] if pd.notna(min_id) else None,
            'value': filtered_df[column].min() if not filtered_df.empty else None,
        }
        max_values[column] = {
            'id': df.loc[max_id, 'id'] if pd.notna(max_id) else None,
            'value': filtered_df[column].max() if not filtered_df.empty else None,
        }

    return min_values, max_values

file_path = 'devices.csv'

min_values, max_values = find_min_max(file_path)

print("Valores mínimos (80% centrais):")
for key, val in min_values.items():
    print(f"{key}: id={val['id']}, valor={val['value']}")

print("\nValores máximos (80% centrais):")
for key, val in max_values.items():
    print(f"{key}: id={val['id']}, valor={val['value']}")
