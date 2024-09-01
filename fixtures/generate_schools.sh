#!/bin/bash

# Fichier de sortie
output_file="/tmp/schools.csv"
size=100

# Générer des noms aléatoires
generate_name() {
  local names=("Springfield Elementary" "Shelbyville Elementary" "Westside High" "Eastside High" "Riverside Academy" "Mountainview School" "Lakeside College" "Hilltop School" "Greenwood High" "Sunset Academy")
  echo "${names[$RANDOM % ${#names[@]}]}"
}

# Générer des adresses aléatoires
generate_address() {
  local streets=("Main St" "Oak St" "Maple Ave" "Pine St" "Cedar Ave" "Birch St" "Elm St" "Willow Ave" "Cherry St" "Peach St")
  local street_number=$((RANDOM % 10000 + 1))
  local street="${streets[$RANDOM % ${#streets[@]}]}"
  echo "$street_number-$street"
}

# Générer les lignes du fichier CSV
for ((i = 1; i <= $size; i++)); do
  name=$(generate_name)
  address=$(generate_address)
  echo "$name;$address" >> "$output_file"
done

echo "Fichier $output_file généré avec succès."