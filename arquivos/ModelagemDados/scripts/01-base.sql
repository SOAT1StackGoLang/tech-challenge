DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'lanchonete') THEN
        CREATE DATABASE lanchonete;
    END IF;
END $$;

\c lanchonete

CREATE TABLE "Cliente" (
    "CodigoCliente" SERIAL,
    "NomeCliente" VARCHAR(100),
    "CPF" VARCHAR(20),
    "Email" VARCHAR(100), 
    CONSTRAINT "PK_Cliente" PRIMARY KEY ("CodigoCliente") 
);

CREATE TABLE "Categoria" (
    "CodigoCategoria" SERIAL,
    "NomeCategoria" VARCHAR(100) NOT NULL,
    CONSTRAINT "PK_Categoria" PRIMARY KEY ("CodigoCategoria") 
);

CREATE TABLE "StatusPedido" (
    "CodigoStatusPedido" SERIAL,
    "NomeStatusPedido" VARCHAR(100) NOT NULL,
    CONSTRAINT "PK_StatusPedido" PRIMARY KEY ("CodigoStatusPedido")
);

CREATE TABLE "StatusPagamento" (
    "CodigoStatusPagamento" SERIAL,
    "NomeStatusPagamento" VARCHAR(100) NOT NULL,
    CONSTRAINT "PK_StatusPagamento" PRIMARY KEY ("CodigoStatusPagamento")
);

CREATE TABLE "Produto" (
    "CodigoProduto" SERIAL,
    "CodigoCategoria" INT NOT NULL,
    "NomeProduto" VARCHAR(100) NOT NULL,
    "DescricaoProduto" VARCHAR(300) NOT NULL,
    "Imagem" VARCHAR(100) NOT NULL,
    "Preco" DECIMAL(8,2) NOT NULL,
    CONSTRAINT "PK_Produto" PRIMARY KEY ("CodigoProduto"),
    CONSTRAINT "FK_Produto_Categoria" FOREIGN KEY ("CodigoCategoria") REFERENCES "Categoria" ("CodigoCategoria") 
);

CREATE TABLE "Pedido" (
    "CodigoPedido" SERIAL,
    "CodigoCliente" INT NOT NULL,
    "DataPedido" TIMESTAMP NOT NULL,
    "DataAtualizacao" TIMESTAMP,
    "CodigoStatusPedido" INT NOT NULL,
    CONSTRAINT "PK_Pedido" PRIMARY KEY ("CodigoPedido"),
    CONSTRAINT "FK_Pedido_Cliente" FOREIGN KEY ("CodigoCliente") REFERENCES "Cliente" ("CodigoCliente"),
    CONSTRAINT "FK_Pedido_StatusPedido" FOREIGN KEY ("CodigoStatusPedido") REFERENCES "StatusPedido" ("CodigoStatusPedido")
);

CREATE TABLE "Pagamento" (
    "CodigoPagamento" SERIAL,
    "CodigoPedido" INT NOT NULL,
    "CodigoStatusPagamento" INT NOT NULL,
    "DataPagamento" TIMESTAMP,
    "Valor" DECIMAL(8,2) NOT NULL,
    "QrCode" VARCHAR(100),
    CONSTRAINT "PK_Pagamento" PRIMARY KEY ("CodigoPagamento"),
    CONSTRAINT "FK_Pagamento_Pedido" FOREIGN KEY ("CodigoPedido") REFERENCES "Pedido" ("CodigoPedido"),
    CONSTRAINT "FK_Pagamento_StatusPagamento" FOREIGN KEY ("CodigoStatusPagamento") REFERENCES "StatusPagamento" ("CodigoStatusPagamento")
);

CREATE TABLE "ItemPedido" (
    "CodigoPedido" INT NOT NULL,
    "CodigoProduto" INT NOT NULL,
    CONSTRAINT "FK_ItemPedido_Pedido" FOREIGN KEY ("CodigoPedido") REFERENCES "Pedido" ("CodigoPedido"),
    CONSTRAINT "FK_ItemPedido_Produto" FOREIGN KEY ("CodigoProduto") REFERENCES "Produto" ("CodigoProduto")
);