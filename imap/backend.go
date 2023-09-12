package imap

// import (
// 	"bufio"
// 	"bytes"
// 	"context"
// 	"fmt"
// 	"io"
// 	"time"
//
// 	"github.com/ItsNotGoodName/smtpbridge/internal/core"
// 	"github.com/ItsNotGoodName/smtpbridge/internal/models"
// 	"github.com/ItsNotGoodName/smtpbridge/pkg/pagination"
// 	imap "github.com/emersion/go-imap"
// 	"github.com/emersion/go-imap/backend"
// )
//
// // The Backend implements SMTP server methods.
// type Backend struct {
// 	app core.App
// }
//
// func NewBackend(app core.App) Backend {
// 	return Backend{
// 		app: app,
// 	}
// }
//
// func (b Backend) Login(connInfo *imap.ConnInfo, username string, password string) (backend.User, error) {
// 	return user{app: b.app}, nil
// }
//
// const defaultMailbox = "INBOX"
//
// var ErrNotAllowed = fmt.Errorf("not allowed")
// var ErrNotImplemented = fmt.Errorf("not implemented")
// var ErrNotFound = fmt.Errorf("not found")
//
// type user struct {
// 	app core.App
// }
//
// // CreateMailbox implements backend.User.
// func (user) CreateMailbox(name string) error {
// 	return ErrNotAllowed
// }
//
// // DeleteMailbox implements backend.User.
// func (user) DeleteMailbox(name string) error {
// 	return ErrNotAllowed
// }
//
// // GetMailbox implements backend.User.
// func (u user) GetMailbox(name string) (backend.Mailbox, error) {
// 	if name != defaultMailbox {
// 		return nil, ErrNotFound
// 	}
//
// 	return mailbox{app: u.app}, nil
// }
//
// // ListMailboxes implements backend.User.
// func (user) ListMailboxes(subscribed bool) ([]backend.Mailbox, error) {
// 	return []backend.Mailbox{mailbox{}}, nil
// }
//
// // Logout implements backend.User.
// func (user) Logout() error {
// 	return nil
// }
//
// // RenameMailbox implements backend.User.
// func (user) RenameMailbox(existingName string, newName string) error {
// 	return ErrNotAllowed
// }
//
// // Username implements backend.User.
// func (user) Username() string {
// 	return ""
// }
//
// type mailbox struct {
// 	app core.App
// }
//
// // Check implements backend.Mailbox.
// func (mailbox) Check() error {
// 	return fmt.Errorf("Check: %w", ErrNotImplemented)
// }
//
// // CopyMessages implements backend.Mailbox.
// func (mailbox) CopyMessages(uid bool, seqset *imap.SeqSet, dest string) error {
// 	return fmt.Errorf("CopyMessages: %w", ErrNotAllowed)
// }
//
// // CreateMessage implements backend.Mailbox.
// func (mailbox) CreateMessage(flags []string, date time.Time, body imap.Literal) error {
// 	return fmt.Errorf("CreateMessage: %w", ErrNotImplemented)
// }
//
// // Expunge implements backend.Mailbox.
// func (mailbox) Expunge() error {
// 	return fmt.Errorf("Expunge: %w", ErrNotImplemented)
// }
//
// // Info implements backend.Mailbox.
// func (mailbox) Info() (*imap.MailboxInfo, error) {
// 	return &imap.MailboxInfo{
// 		Name: defaultMailbox,
// 	}, nil
// }
//
// // ListMessages implements backend.Mailbox.
// func (m mailbox) ListMessages(uid bool, seqSet *imap.SeqSet, items []imap.FetchItem, ch chan<- *imap.Message) error {
// 	defer close(ch)
// 	ctx := context.Background()
//
// 	list, err := m.app.EnvelopeList(ctx, pagination.NewPage(1, 100), models.DTOEnvelopeListRequest{})
// 	if err != nil {
// 		return err
// 	}
//
// 	for i, msg := range list.Envelopes {
// 		seqNum := uint32(i + 1)
//
// 		var id uint32
// 		if uid {
// 			id = uint32(msg.Message.ID)
// 		} else {
// 			id = seqNum
// 		}
// 		if !seqSet.Contains(id) {
// 			continue
// 		}
//
// 		m, err := fetch(msg, seqNum, items)
// 		if err != nil {
// 			continue
// 		}
//
// 		ch <- m
// 	}
//
// 	return nil
// }
//
// // Name implements backend.Mailbox.
// func (mailbox) Name() string {
// 	return defaultMailbox
// }
//
// // SearchMessages implements backend.Mailbox.
// func (mailbox) SearchMessages(uid bool, criteria *imap.SearchCriteria) ([]uint32, error) {
// 	return nil, fmt.Errorf("SearchMessages: %w", ErrNotImplemented)
// }
//
// // SetSubscribed implements backend.Mailbox.
// func (mailbox) SetSubscribed(subscribed bool) error {
// 	return fmt.Errorf("SetSubscribed: %w", ErrNotImplemented)
// }
//
// // Status implements backend.Mailbox.
// func (m mailbox) Status(items []imap.StatusItem) (*imap.MailboxStatus, error) {
// 	ctx := context.Background()
// 	status := imap.NewMailboxStatus(defaultMailbox, items)
//
// 	for _, name := range items {
// 		switch name {
// 		case imap.StatusMessages:
// 			storage, err := m.app.StorageGet(ctx)
// 			if err != nil {
// 				return nil, err
// 			}
// 			status.Messages = uint32(storage.EnvelopeCount)
// 		case imap.StatusUidNext:
// 			env, err := m.app.EnvelopeList(ctx, pagination.NewPage(1, 1), models.DTOEnvelopeListRequest{})
// 			if err != nil {
// 				return nil, err
// 			}
//
// 			var uid uint32 = 1
// 			if len(env.Envelopes) > 0 {
// 				uid = uint32(env.Envelopes[0].Message.ID)
// 			}
//
// 			status.UidNext = uid
// 		case imap.StatusUidValidity:
// 			status.UidValidity = 1
// 		case imap.StatusRecent:
// 			status.Recent = 0 // TODO
// 		case imap.StatusUnseen:
// 			status.Unseen = 0 // TODO
// 		}
// 	}
//
// 	return status, nil
// }
//
// // UpdateMessagesFlags implements backend.Mailbox.
// func (mailbox) UpdateMessagesFlags(uid bool, seqset *imap.SeqSet, operation imap.FlagsOp, flags []string) error {
// 	return fmt.Errorf("UpdateMessagesFlags: %w", ErrNotAllowed)
// }
//
// func fetch(env models.Envelope, seqNum uint32, items []imap.FetchItem) (*imap.Message, error) {
// 	fetched := imap.NewMessage(seqNum, items)
// 	for _, item := range items {
// 		switch item {
// 		case imap.FetchEnvelope:
// 			var to []*imap.Address
// 			for _, v := range env.Message.To {
// 				to = append(to, &imap.Address{
// 					MailboxName:  defaultMailbox,
// 					PersonalName: v,
// 				})
// 			}
//
// 			fetched.Envelope = &imap.Envelope{
// 				Date:    env.Message.Date.Time(),
// 				Subject: env.Message.Subject,
// 				To:      to,
// 				From:    []*imap.Address{{PersonalName: env.Message.From}},
// 			}
// 		case imap.FetchBody, imap.FetchBodyStructure:
// 			// hdr, body, _ := m.headerAndBody()
// 			// fetched.BodyStructure, _ = backendutil.FetchBodyStructure(hdr, body, item == imap.FetchBodyStructure)
//
// 			if item == imap.FetchBodyStructure {
// 				panic("nope")
// 			}
//
// 			fetched.BodyStructure = &imap.BodyStructure{
// 				MIMEType: "text",
// 			}
// 		case imap.FetchFlags:
// 			fetched.Flags = []string{}
// 		case imap.FetchInternalDate:
// 			fetched.InternalDate = env.Message.CreatedAt.Time()
// 		case imap.FetchRFC822Size:
// 			fetched.Size = 0 // TODO:
// 		case imap.FetchUid:
// 			fetched.Uid = uint32(env.Message.ID)
// 		default:
// 			section, err := imap.ParseBodySectionName(item)
// 			if err != nil {
// 				break
// 			}
//
// 			// hdr, err := textproto.ReadHeader(body)
// 			// if err != nil {
// 			// 	return nil, err
// 			// }
// 			//
// 			// l, _ := backendutil.FetchBodySection(hdr, body, section)
//
// 			fetched.Body[section] = newMessageBody(env.Message.Text)
// 		}
// 	}
//
// 	return fetched, nil
// }
//
// type messagebody struct {
// 	io.Reader
// 	length int
// }
//
// func newMessageBody(t string) messagebody {
// 	length := len(t)
// 	body := bufio.NewReader(bytes.NewReader([]byte(t)))
// 	return messagebody{
// 		Reader: body,
// 		length: length,
// 	}
// }
//
// func (m messagebody) Len() int {
// 	return m.length
// }
//
// //
// // // Append implements imapserver.Session.
// // func (session) Append(mailbox string, r imap.LiteralReader, options *imap.AppendOptions) (*imap.AppendData, error) {
// // 	return nil, ErrNotAllowed
// // }
// //
// // // Close implements imapserver.Session.
// // func (session) Close() error {
// // 	return nil
// // }
// //
// // // Copy implements imapserver.Session.
// // func (session) Copy(kind imapserver.NumKind, seqSet imap.SeqSet, dest string) (*imap.CopyData, error) {
// // 	return nil, ErrNotAllowed
// // }
// //
// // // Create implements imapserver.Session.
// // func (session) Create(mailbox string, options *imap.CreateOptions) error {
// // 	return ErrNotAllowed
// // }
// //
// // // Delete implements imapserver.Session.
// // func (session) Delete(mailbox string) error {
// // 	return ErrNotAllowed
// // }
// //
// // // Expunge implements imapserver.Session.
// // func (session) Expunge(w *imapserver.ExpungeWriter, uids *imap.SeqSet) error {
// // 	return ErrNotImplemented
// // }
// //
// // // Fetch implements imapserver.Session.
// // func (s session) Fetch(w *imapserver.FetchWriter, kind imapserver.NumKind, seqSet imap.SeqSet, options *imap.FetchOptions) error {
// //
// // 	fmt.Println(helpers.JSON(kind))
// //
// // 	ctx := context.Background()
// // 	for _, seq := range seqSet {
// // 		env, err := s.app.EnvelopeGet(ctx, int64(seq.Start))
// // 		if err != nil {
// // 			return err
// // 		}
// //
// // 		wt := w.CreateMessage(seq.Start)
// // 		defer wt.Close()
// //
// // 		var to []imap.Address
// // 		for _, v := range env.Message.To {
// // 			to = append(to, imap.Address{
// // 				Mailbox: defaultMailbox,
// // 				Name:    v,
// // 			})
// // 		}
// //
// // 		wt.WriteEnvelope(&imap.Envelope{
// // 			Subject: env.Message.Subject,
// // 			From:    []imap.Address{{Name: env.Message.From, Mailbox: defaultMailbox}},
// // 			To:      to,
// // 			Date:    env.Message.Date.Time(),
// // 		})
// //
// // 		wt.Close()
// // 	}
// //
// // 	return nil
// // }
// //
// // // Idle implements imapserver.Session.
// // func (session) Idle(w *imapserver.UpdateWriter, stop <-chan struct{}) error {
// // 	<-stop
// // 	return nil
// // }
// //
// // // List implements imapserver.Session.
// // func (session) List(w *imapserver.ListWriter, ref string, patterns []string, options *imap.ListOptions) error {
// // 	return ErrNotImplemented
// // }
// //
// // // Login implements imapserver.Session.
// // func (session) Login(username string, password string) error {
// // 	return nil
// // }
// //
// // // Poll implements imapserver.Session.
// // func (session) Poll(w *imapserver.UpdateWriter, allowExpunge bool) error {
// // 	return ErrNotAllowed
// // }
// //
// // // Rename implements imapserver.Session.
// // func (session) Rename(mailbox string, newName string) error {
// // 	return ErrNotAllowed
// // }
// //
// // // Search implements imapserver.Session.
// // func (session) Search(kind imapserver.NumKind, criteria *imap.SearchCriteria, options *imap.SearchOptions) (*imap.SearchData, error) {
// // 	return nil, ErrNotImplemented
// // }
// //
// // // Select implements imapserver.Session.
// // func (s session) Select(mailbox string, options *imap.SelectOptions) (*imap.SelectData, error) {
// // 	if mailbox != defaultMailbox {
// // 		return nil, ErrNotFound
// // 	}
// //
// // 	env, err := s.app.EnvelopeList(context.Background(), pagination.NewPage(1, 100), models.DTOEnvelopeListRequest{})
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	var uidNext uint32
// // 	if len(env.Envelopes) > 0 {
// // 		uidNext = uint32(env.Envelopes[0].Message.ID)
// // 	}
// //
// // 	return &imap.SelectData{
// // 		NumMessages: uint32(env.PageResult.TotalItems),
// // 		UIDValidity: uidNext,
// // 		UIDNext:     uidNext,
// // 	}, nil
// // }
// //
// // // Status implements imapserver.Session.
// // func (session) Status(mailbox string, options *imap.StatusOptions) (*imap.StatusData, error) {
// // 	return nil, ErrNotImplemented
// // }
// //
// // // Store implements imapserver.Session.
// // func (session) Store(w *imapserver.FetchWriter, kind imapserver.NumKind, seqSet imap.SeqSet, flags *imap.StoreFlags, options *imap.StoreOptions) error {
// // 	return ErrNotImplemented
// // }
// //
// // // Subscribe implements imapserver.Session.
// // func (session) Subscribe(mailbox string) error {
// // 	return ErrNotAllowed
// // }
// //
// // // Unselect implements imapserver.Session.
// // func (session) Unselect() error {
// // 	return ErrNotImplemented
// // }
// //
// // // Unsubscribe implements imapserver.Session.
// // func (session) Unsubscribe(mailbox string) error {
// // 	return ErrNotAllowed
// // }
// //
// // var _ imapserver.Session = session{}
