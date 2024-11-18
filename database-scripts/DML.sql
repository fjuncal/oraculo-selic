CREATE TABLE MENSAGENS (
                           id SERIAL PRIMARY KEY,
                           TXT_CORREL_ID UUID DEFAULT uuid_generate_v4(),   -- Identificador único, gerado automaticamente na criação
                           TXT_COD_MSG VARCHAR(50) NOT NULL,
                           TXT_MSG_DOC_XML TEXT,
                           TXT_MSG TEXT,
                           TXT_CANAL TEXT,
                           TXT_STATUS VARCHAR(50),
                           DT_INCL TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE PASSOS_TESTES (
                          id SERIAL PRIMARY KEY,
                          TXT_DESCRICAO TEXT,
                          TXT_TP_PASSO_TESTE TEXT,
                          TXT_CANAL TEXT,
                          TXT_COD_MSG VARCHAR(10) NOT NULL,                -- Código da mensagem, como SEL1052
                          TXT_MSG_DOC_XML TEXT,                          -- Conteúdo completo em XML
                          TXT_MSG TEXT,                          -- String completa no formato SELIC
                          TXT_CT_CED TEXT,
                          TXT_CT_CESS TEXT,
                          TXT_NUM_OP TEXT,
                          TXT_EMISSOR TEXT,
                          VAL_FIN NUMERIC(10, 2),            -- Valor financeiro, caso aplicável
                          VAL_PU NUMERIC(10, 2),    -- Valor financeiro de retorno, se aplicável
                          DT_INCL TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE CENARIOS (
                          id SERIAL PRIMARY KEY,
                          TXT_NOME VARCHAR(255) NOT NULL, -- Nome do cenário
                          TXT_DESCRICAO TEXT,             -- Descrição do cenário
                          DT_INCL TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE CENARIOS_PASSOS_TESTES (
                                        ID_CENARIO INTEGER NOT NULL REFERENCES CENARIOS(id) ON DELETE CASCADE,
                                        ID_PASSO_TESTE INTEGER NOT NULL REFERENCES PASSOS_TESTES(id) ON DELETE CASCADE,
                                        ORDENACAO INTEGER,             -- Ordem do passo dentro do cenário
                                        DT_INCL TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                        PRIMARY KEY (ID_CENARIO, ID_PASSO_TESTE)
);