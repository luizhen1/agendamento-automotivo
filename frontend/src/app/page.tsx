"use client";

import { FormEvent, useCallback, useEffect, useState } from "react";
import styles from "./page.module.css";

const API = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";
type Cliente = { id: number; nome: string; telefone: string; veiculo_modelo: string; placa: string };
type Servico = { id: number; nome: string; descricao: string; valor: number; duracao_minutos: number };
type Agendamento = { id: number; cliente_id: number; servico_id: number; data_hora: string; status: string; observacoes: string; cliente: Cliente; servico: Servico };
type Feedback = { kind: "success" | "error"; message: string };

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API}${path}`, { ...init, headers: { "Content-Type": "application/json", ...init?.headers } });
  const body = await response.json().catch(() => ({}));
  if (!response.ok) throw new Error(body.erro || `Erro HTTP ${response.status}`);
  return body as T;
}
const money = (value: number) => new Intl.NumberFormat("pt-BR", { style: "currency", currency: "BRL" }).format(value);
const date = (value: string) => {
  const parsed = new Date(value);
  return Number.isNaN(parsed.getTime()) ? value.replace("T", " ") : parsed.toLocaleString("pt-BR", { dateStyle: "short", timeStyle: "short" });
};

export default function Home() {
  const [clientes, setClientes] = useState<Cliente[]>([]);
  const [servicos, setServicos] = useState<Servico[]>([]);
  const [agendamentos, setAgendamentos] = useState<Agendamento[]>([]);
  const [feedback, setFeedback] = useState<Feedback | null>(null);
  const [busy, setBusy] = useState(false);

  const load = useCallback(async () => {
    try {
      const [c, s, a] = await Promise.all([request<Cliente[]>("/clientes"), request<Servico[]>("/servicos"), request<Agendamento[]>("/agendamentos")]);
      setClientes(c); setServicos(s); setAgendamentos(a);
    } catch (error) { setFeedback({ kind: "error", message: error instanceof Error ? error.message : "Não foi possível carregar os dados." }); }
  }, []);
  useEffect(() => { void load(); }, [load]);

  async function submit(event: FormEvent<HTMLFormElement>, path: string, success: string, transform?: (data: Record<string, FormDataEntryValue>) => object) {
    event.preventDefault(); setBusy(true); setFeedback(null);
    const formElement = event.currentTarget;
    const form = new FormData(formElement);
    const values = Object.fromEntries(form.entries());
    try {
      await request(path, { method: "POST", body: JSON.stringify(transform ? transform(values) : values) });
      formElement.reset(); setFeedback({ kind: "success", message: success }); await load();
    } catch (error) { setFeedback({ kind: "error", message: error instanceof Error ? error.message : "Não foi possível concluir a operação." }); }
    finally { setBusy(false); }
  }

  async function update(id: number, payload: object, success: string) {
    setFeedback(null);
    try { await request(`/agendamentos/${id}`, { method: "PATCH", body: JSON.stringify(payload) }); setFeedback({ kind: "success", message: success }); await load(); }
    catch (error) { setFeedback({ kind: "error", message: error instanceof Error ? error.message : "Falha ao atualizar agendamento." }); }
  }
  async function remove(id: number) {
    if (!window.confirm("Deseja excluir este agendamento?")) return;
    setFeedback(null);
    try { await request(`/agendamentos/${id}`, { method: "DELETE" }); setFeedback({ kind: "success", message: "Agendamento excluído." }); await load(); }
    catch (error) { setFeedback({ kind: "error", message: error instanceof Error ? error.message : "Falha ao excluir agendamento." }); }
  }

  return <main className={styles.shell}>
    <header className={styles.header}>
      <div><div className={styles.eyebrow}>CENTRO AUTOMOTIVO · PAINEL DE GESTÃO</div><h1>Sistema de Agendamento de Serviços Automotivos</h1><p>Gestão simples para sua oficina, do cadastro à entrega do veículo.</p></div>
      <div className={styles.stack}>FULL-STACK<br /><strong>Go · MySQL · Next.js</strong></div>
    </header>

    {feedback && <div role="alert" className={`${styles.feedback} ${feedback.kind === "success" ? styles.success : styles.error}`}><span>{feedback.kind === "success" ? "✓" : "!"}</span>{feedback.message}<button aria-label="Fechar aviso" onClick={() => setFeedback(null)}>×</button></div>}

    <section className={styles.forms}>
      <article className={styles.card}>
        <div className={styles.cardHeading}><span className={styles.number}>01</span><div><h2>Cadastrar novo serviço</h2><p>UC01 · Adicione uma opção ao catálogo</p></div></div>
        <form className={styles.form} onSubmit={(e) => submit(e, "/servicos", "Serviço cadastrado com sucesso!", (v) => ({ nome: v.nome, descricao: v.descricao, valor: Number(v.valor), duracao_minutos: Number(v.duracao_minutos) }))}>
          <label>Nome do serviço<input name="nome" required maxLength={100} placeholder="Ex.: Lavagem detalhada" /></label>
          <label>Descrição<textarea name="descricao" maxLength={255} rows={2} placeholder="Descreva o que está incluso" /></label>
          <div className={styles.twoFields}><label>Valor (R$)<input name="valor" type="number" min="0.01" step="0.01" required placeholder="0,00" /></label><label>Duração (minutos)<input name="duracao_minutos" type="number" min="1" step="1" required placeholder="60" /></label></div>
          <button className={styles.primary} disabled={busy}>Cadastrar serviço <span>→</span></button>
        </form>
      </article>

      <article className={styles.card}>
        <div className={styles.cardHeading}><span className={`${styles.number} ${styles.numberViolet}`}>03</span><div><h2>Realizar agendamento</h2><p>UC03 · Reserve um horário para o cliente</p></div></div>
        <form className={styles.form} onSubmit={(e) => submit(e, "/agendamentos", "Agendamento realizado com sucesso!", (v) => ({ cliente_id: Number(v.cliente_id), servico_id: Number(v.servico_id), data_hora: v.data_hora, observacoes: v.observacoes }))}>
          <label>Cliente<select name="cliente_id" required defaultValue=""><option value="" disabled>Selecione um cliente</option>{clientes.map(c => <option key={c.id} value={c.id}>{c.nome} · {c.veiculo_modelo} · {c.placa}</option>)}</select></label>
          <label>Serviço<select name="servico_id" required defaultValue=""><option value="" disabled>Selecione um serviço</option>{servicos.map(s => <option key={s.id} value={s.id}>{s.nome} · {money(s.valor)} · {s.duracao_minutos} min</option>)}</select></label>
          <div className={styles.twoFields}><label>Data e hora<input name="data_hora" type="datetime-local" required /></label><label>Observações<input name="observacoes" maxLength={255} placeholder="Opcional" /></label></div>
          <button className={`${styles.primary} ${styles.violet}`} disabled={busy}>Confirmar agendamento <span>→</span></button>
        </form>
      </article>
    </section>

    <section className={styles.section}>
      <div className={styles.sectionHeading}><div><div className={styles.eyebrow}>UC02 · UC04 · UC05</div><h2>Agenda de serviços</h2><p>Acompanhe os próximos atendimentos e atualize cada serviço.</p></div><div className={styles.count}>{agendamentos.length} {agendamentos.length === 1 ? "agendamento" : "agendamentos"}</div></div>
      <div className={styles.tableWrap}><table><thead><tr><th>ID</th><th>Data e hora</th><th>Cliente / veículo</th><th>Serviço</th><th>Valor</th><th>Status</th><th>Ações</th></tr></thead><tbody>
        {agendamentos.map(a => <tr key={a.id}><td className={styles.id}>#{a.id}</td><td>{date(a.data_hora)}</td><td><strong>{a.cliente?.nome ?? "—"}</strong><small>{a.cliente?.veiculo_modelo} · {a.cliente?.placa}</small></td><td>{a.servico?.nome ?? "—"}</td><td>{money(a.servico?.valor ?? 0)}</td><td><span className={`${styles.badge} ${a.status === "CONCLUIDO" ? styles.done : a.status === "CANCELADO" ? styles.cancelled : styles.scheduled}`}><i />{a.status}</span></td><td><div className={styles.actions}>{a.status === "AGENDADO" && <button className={styles.complete} onClick={() => update(a.id, { status: "CONCLUIDO" }, "Agendamento concluído.")}>Concluir (UC04)</button>}<button className={styles.delete} onClick={() => remove(a.id)}>Excluir (UC05)</button></div></td></tr>)}
        {agendamentos.length === 0 && <tr><td colSpan={7} className={styles.empty}>Nenhum agendamento cadastrado. Use o formulário acima para começar.</td></tr>}
      </tbody></table></div>
    </section>

    <section className={styles.section}>
      <div className={styles.sectionHeading}><div><div className={styles.eyebrow}>UC02 · CATÁLOGO</div><h2>Serviços disponíveis</h2><p>Opções cadastradas para atendimento.</p></div><div className={styles.count}>{servicos.length} serviços</div></div>
      <div className={styles.tableWrap}><table><thead><tr><th>ID</th><th>Nome</th><th>Descrição</th><th>Duração</th><th>Valor</th></tr></thead><tbody>
        {servicos.map(s => <tr key={s.id}><td className={styles.id}>#{s.id}</td><td><strong>{s.nome}</strong></td><td>{s.descricao || "—"}</td><td>{s.duracao_minutos} min</td><td className={styles.price}>{money(s.valor)}</td></tr>)}
        {servicos.length === 0 && <tr><td colSpan={5} className={styles.empty}>Nenhum serviço disponível.</td></tr>}
      </tbody></table></div>
    </section>
    <footer className={styles.footer}>Sistema de Agendamento Automotivo <span>·</span> Go + GORM + MySQL <span>·</span> Next.js</footer>
  </main>;
}
