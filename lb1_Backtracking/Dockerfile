FROM mcr.microsoft.com/dotnet/runtime:9.0 AS base
WORKDIR /app

FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build
WORKDIR /src

COPY ["lb1_Backtracking.csproj", "./"]
RUN dotnet restore "lb1_Backtracking.csproj"

COPY . .
RUN dotnet build "lb1_Backtracking.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "lb1_Backtracking.csproj" -c Release -o /app/publish

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .

CMD ["bash"]