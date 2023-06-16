CREATE TABLE IF NOT EXISTS Cliente (
    CodigoCliente SERIAL PRIMARY KEY,
    NomeCliente VARCHAR(100),
    CPF VARCHAR(20),
    Email VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS Produto (
    CodigoProduto SERIAL PRIMARY KEY,
    CodigoCategoria INTEGER,
    NomeProduto VARCHAR(100),
    DescricaoProduto VARCHAR(300),
    Imagem VARCHAR(100),
    Preco DECIMAL
);

CREATE TABLE IF NOT EXISTS Categoria (
    CodigoCategoria SERIAL PRIMARY KEY,
    NomeCategoria VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS Pedido (
    CodigoPedido SERIAL PRIMARY KEY,
    CodigoCliente INTEGER,
    DataPedido TIMESTAMP,
    DataAtualizacao TIMESTAMP,
    CodigoStatusPedido INTEGER
);

CREATE TABLE IF NOT EXISTS Pagamento (
    CodigoPagamento SERIAL PRIMARY KEY,
    CodigoPedido INTEGER,
    CodigoStatusPagamento INTEGER,
    DataPagamento TIMESTAMP,
    QrCode VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS Status_Pedido (
    CodigoStatusPedido SERIAL PRIMARY KEY,
    NomeStatusPedido VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS Status_Pagamento (
    CodigoStatusPagamento SERIAL PRIMARY KEY,
    NomeStatusPagamento VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS ItemPedido (
    CodigoPedido INTEGER,
    CodigoProduto INTEGER
);
 
ALTER TABLE Produto ADD CONSTRAINT FK_Produto_2
    FOREIGN KEY (CodigoCategoria)
    REFERENCES Categoria (CodigoCategoria)
    ON DELETE CASCADE;
 
ALTER TABLE Pedido ADD CONSTRAINT FK_Pedido_2
    FOREIGN KEY (CodigoStatusPedido)
    REFERENCES Status_Pedido (CodigoStatusPedido)
    ON DELETE CASCADE;
 
ALTER TABLE Pedido ADD CONSTRAINT FK_Pedido_3
    FOREIGN KEY (CodigoCliente)
    REFERENCES Cliente (CodigoCliente)
    ON DELETE CASCADE;
 
ALTER TABLE Pagamento ADD CONSTRAINT FK_Pagamento_2
    FOREIGN KEY (CodigoStatusPagamento)
    REFERENCES Status_Pagamento (CodigoStatusPagamento)
    ON DELETE CASCADE;
 
ALTER TABLE Pagamento ADD CONSTRAINT FK_Pagamento_3
    FOREIGN KEY (CodigoPedido)
    REFERENCES Pedido (CodigoPedido)
    ON DELETE CASCADE;
 
ALTER TABLE ItemPedido ADD CONSTRAINT FK_ItemPedido_1
    FOREIGN KEY (CodigoProduto)
    REFERENCES Produto (CodigoProduto)
    ON DELETE SET NULL;
 
ALTER TABLE ItemPedido ADD CONSTRAINT FK_ItemPedido_2
    FOREIGN KEY (CodigoPedido)
    REFERENCES Pedido (CodigoPedido)
    ON DELETE SET NULL;